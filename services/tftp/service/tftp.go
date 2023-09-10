package service

import (
	"context"
	"errors"
	"io"
	"net"
	stdhttp "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pin/tftp/v3"
	"github.com/shurcooL/githubv4"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/z5labs/app"
	"github.com/z5labs/app/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

type Duration time.Duration

type config struct {
	Tftp struct {
		Port      uint          `config:"port"`
		BlockSize int           `config:"blockSize"`
		Retries   int           `config:"retries"`
		Timeout   time.Duration `config:"timeout"`
	} `config:"tftp"`

	Github struct {
		Token string `config:"token"`
	} `config:"github"`
}

type runtime struct {
	log *otelzap.Logger

	// tftp
	port      uint
	blockSize int
	retries   int
	timeout   time.Duration

	http   *stdhttp.Client
	github *githubv4.Client
}

func BuildRuntime(bc app.BuildContext) (app.Runtime, error) {
	var cfg config
	err := bc.Config.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Github.Token},
	)
	oauthTransport := &oauth2.Transport{
		Source: src,
		Base:   stdhttp.DefaultTransport,
	}
	hc := http.NewClient(
		http.WithTransport(
			http.RoundTripperWith(
				oauthTransport,
				http.CircuitBreaker(
					http.CircuitLogger(logger),
				),
			),
		),
		http.RetryRequests(
			http.MinWaitDuration(100*time.Millisecond),
			http.MaxWaitDuration(cfg.Tftp.Timeout),
			http.MaxAttempts(1),
		),
	)

	rt := &runtime{
		log:       otelzap.New(logger),
		port:      cfg.Tftp.Port,
		blockSize: cfg.Tftp.BlockSize,
		retries:   cfg.Tftp.Retries,
		timeout:   cfg.Tftp.Timeout,
		http:      hc,
		github:    githubv4.NewClient(hc),
	}
	return rt, nil
}

func (rt *runtime) Run(ctx context.Context) error {
	conn, err := net.ListenPacket("udp", ":"+strconv.FormatUint(uint64(rt.port), 10))
	if err != nil {
		return err
	}

	s := tftp.NewServer(rt.readHandler, nil)
	s.SetAnticipate(2)
	s.SetBlockSize(rt.blockSize)
	s.SetRetries(rt.retries)
	s.SetTimeout(rt.timeout)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gctx.Done()
		s.Shutdown()
		return nil
	})
	g.Go(func() error {
		rt.log.Info("started serving tftp")
		defer rt.log.Info("stopped service tftp")
		return s.Serve(conn)
	})
	return g.Wait()
}

func (rt *runtime) readHandler(filename string, rf io.ReaderFrom) error {
	spanCtx, span := otel.Tracer("service").Start(context.Background(), "runtime.readHandler", trace.WithAttributes(
		attribute.String("filename", filename),
	))
	defer span.End()

	ss := strings.Split(filename, "/")
	if len(ss) < 4 {
		return errors.New("filepath not long enough")
	}
	ss = ss[len(ss)-4:] // TODO: sanitize these values

	log := rt.log.Ctx(spanCtx).WithOptions(zap.Fields(
		zap.String("github_repo_owner", ss[0]),
		zap.String("github_repo_name", ss[1]),
		zap.String("github_release_tag", ss[2]),
		zap.String("github_release_asset_name", ss[3]),
	))

	log.Info("querying github for release asset download url")
	var query struct {
		Repository struct {
			Release struct {
				ID            githubv4.ID       `graphql:"id"`
				CreatedAt     githubv4.DateTime `graphql:"createdAt"`
				TagName       string            `graphql:"tagName"`
				ReleaseAssets struct {
					Nodes []struct {
						ID          githubv4.ID       `graphql:"id"`
						CreatedAt   githubv4.DateTime `graphql:"createdAt"`
						Size        int               `graphql:"size"`
						ContentType string            `graphql:"contentType"`
						DownloadUrl string            `graphql:"downloadUrl"`
					} `graphql:"nodes"`
				} `graphql:"releaseAssets(first: 1, name: $assetName)"`
			} `graphql:"release(tagName: $releaseTag)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}
	err := rt.github.Query(spanCtx, &query, map[string]interface{}{
		"$repoOwner":  ss[0],
		"$repoName":   ss[1],
		"$releaseTag": ss[2],
		"$assetName":  ss[3],
	})
	if err != nil {
		log.Error("failed to query github for release asset download url", zap.Error(err))
		return errors.New("query failed")
	}

	assets := query.Repository.Release.ReleaseAssets.Nodes
	if len(assets) == 0 {
		log.Error("no matching release asset found")
		return errors.New("no release asset found")
	}
	asset := assets[0]

	req, err := stdhttp.NewRequest(stdhttp.MethodGet, asset.DownloadUrl, nil)
	if err != nil {
		log.Error("failed to build http request", zap.Error(err))
		return err
	}
	req = req.WithContext(spanCtx)

	resp, err := rt.http.Do(req)
	if err != nil {
		log.Error("failed to perform http request", zap.Error(err))
		return err
	}
	if resp.StatusCode != stdhttp.StatusOK {
		log.Error("received unexpected http status code from github download endpoint", zap.Int("http_status_code", resp.StatusCode))
		return errors.New("unexpected http status code")
	}
	defer resp.Body.Close()

	n, err := rf.ReadFrom(resp.Body)
	if err != nil {
		log.Error("failed to read from response body", zap.Error(err))
		return err
	}
	if n != int64(asset.Size) {
		log.Error("failed to read all of the response body", zap.Int64("bytes_read", n), zap.Int("response_body_bytes", asset.Size))
		return errors.New("did not read all of response body")
	}
	return nil
}
