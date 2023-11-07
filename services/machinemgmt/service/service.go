package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/Zaba505/infra/pkg/httpvalidate"
	"github.com/Zaba505/infra/services/machinemgmt/service/backend"

	"cloud.google.com/go/storage"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/z5labs/app"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type config struct {
	Http struct {
		Port uint `config:"port"`
	} `config:"http"`

	Storage struct {
		Bucket string `config:"bucket"`
	} `config:"storage"`
}

type storageClient interface {
	GetBootstrapImage(context.Context, *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error)
}

type runtime struct {
	log    *otelzap.Logger
	listen func(string, string) (net.Listener, error)

	// http
	port uint

	started atomic.Bool
	healthy atomic.Bool
	serving atomic.Bool

	storage storageClient
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

	tp, err := configureOtel()
	if err != nil {
		logger.Error("failed to configure otel", zap.Error(err))
		return nil, err
	}
	otel.SetTracerProvider(tp)

	gs, err := storage.NewClient(context.Background())
	if err != nil {
		logger.Error("failed to create storage client", zap.Error(err))
		return nil, err
	}
	bucket := gs.Bucket(cfg.Storage.Bucket)
	storageService := backend.NewStorageService(
		backend.Logger(logger),
		backend.GoogleCloudBucket(bucket),
		backend.ObjectHasher(sha256.New),
	)

	rt := &runtime{
		log:     otelzap.New(logger),
		listen:  net.Listen,
		port:    cfg.Http.Port,
		storage: storageService,
	}
	return rt, nil
}

func (rt *runtime) Run(ctx context.Context) error {
	conn, err := rt.listen("tcp", ":"+strconv.FormatUint(uint64(rt.port), 10))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	registerEndpoint(
		mux,
		"/health/startup",
		httpvalidate.Request(
			http.HandlerFunc(rt.startupHandler),
			httpvalidate.ForMethods(http.MethodGet),
		),
	)
	registerEndpoint(
		mux,
		"/health/liveness",
		httpvalidate.Request(
			http.HandlerFunc(rt.livenessHandler),
			httpvalidate.ForMethods(http.MethodGet),
		),
	)
	registerEndpoint(
		mux,
		"/health/readiness",
		httpvalidate.Request(
			http.HandlerFunc(rt.readinessHandler),
			httpvalidate.ForMethods(http.MethodGet),
		),
	)
	registerEndpoint(
		mux,
		"/bootstrap/image",
		httpvalidate.Request(
			http.HandlerFunc(rt.bootstrapImageHandler),
			httpvalidate.ForMethods(http.MethodGet),
			httpvalidate.ExactParams("id"),
		),
	)

	s := &http.Server{
		Handler: otelhttp.NewHandler(mux, "machinemgmt", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gctx.Done()
		defer func() {
			tp := otel.GetTracerProvider()
			stp, ok := tp.(interface {
				Shutdown(context.Context) error
			})
			if !ok {
				return
			}
			stp.Shutdown(context.Background())
		}()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer rt.log.Info("stopped service")

		s.Shutdown(ctx)
		return nil
	})
	g.Go(func() error {
		rt.log.Info("started serving")
		rt.started.Store(true)
		rt.healthy.Store(true)
		rt.serving.Store(true)
		return s.Serve(conn)
	})

	err = g.Wait()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func registerEndpoint(mux *http.ServeMux, path string, h http.Handler) {
	mux.Handle(
		path,
		otelhttp.WithRouteTag(path, h),
	)
}

// report whether this service is ready to begin accepting traffic
func (rt *runtime) startupHandler(w http.ResponseWriter, req *http.Request) {
	started := rt.started.Load()
	if started {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// report whether this service is healthy or needs to be restarted
func (rt *runtime) livenessHandler(w http.ResponseWriter, req *http.Request) {
	healthy := rt.healthy.Load()
	if healthy {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// report whether this service is able to accept traffic
func (rt *runtime) readinessHandler(w http.ResponseWriter, req *http.Request) {
	serving := rt.serving.Load()
	if serving {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (rt *runtime) bootstrapImageHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	imageId := params.Get("id")

	spanCtx, span := otel.Tracer("service").Start(req.Context(), "runtime.bootstrapImageHandler", trace.WithAttributes(
		attribute.String("image.id", imageId),
	))
	defer span.End()

	resp, err := rt.storage.GetBootstrapImage(spanCtx, &backend.GetBootstrapImageRequest{
		ID: imageId,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rt.log.Ctx(spanCtx).Error("failed to get bootstrap image", zap.String("image_id", imageId), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	base64Hash := base64.URLEncoding.EncodeToString(resp.Hash)
	w.Header().Add("ETag", fmt.Sprintf("sha256/%s", base64Hash))
	w.Header().Add("Content-Type", "application/octet")

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rt.log.Ctx(spanCtx).Error("failed to write image to response", zap.String("image_id", imageId), zap.Error(err))
		return
	}
}
