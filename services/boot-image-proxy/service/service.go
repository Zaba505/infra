package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/Zaba505/infra/pkg/framework"
	"github.com/z5labs/bedrock/http/httpclient"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"go.opentelemetry.io/otel"
	"goftp.io/server/v2"
)

type Config struct {
	framework.Config `config:",squash"`

	Proxy struct {
		Target string `config:"target"`
	} `config:"proxy"`
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type httpProxyDriver struct {
	log    *slog.Logger
	http   httpClient
	target string
}

func Init(ctx context.Context) (server.Driver, error) {
	logHandler := framework.LogHandler()
	log := slog.New(logHandler)

	var cfg Config
	err := framework.UnmarshalConfigFromContext(ctx, &cfg)
	if err != nil {
		log.ErrorContext(ctx, "failed to unmarshal config", slogfield.Error(err))
		return nil, err
	}

	d := &httpProxyDriver{
		log:    log,
		http:   httpclient.New(),
		target: cfg.Proxy.Target,
	}
	return d, nil
}

// GetFile implements server.Driver.
func (d *httpProxyDriver) GetFile(ctx *server.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	spanCtx, span := otel.Tracer("service").Start(context.Background(), "httpProxyDriver.GetFile")
	defer span.End()

	req, err := http.NewRequestWithContext(spanCtx, http.MethodGet, fmt.Sprintf("https://%s/%s", d.target, path), nil)
	if err != nil {
		d.log.ErrorContext(spanCtx, "failed to construct http request", slogfield.Error(err))
		return 0, nil, err
	}

	resp, err := d.http.Do(req)
	if err != nil {
		d.log.ErrorContext(spanCtx, "http request failed", slogfield.Error(err))
		return 0, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d.log.ErrorContext(
			spanCtx,
			"received unexpected http status code from backend",
			slogfield.Int("http_status_code", resp.StatusCode),
		)
		return 0, nil, errors.New("unexpected http status code")
	}

	return resp.ContentLength, resp.Body, nil
}

var errUnsupported = errors.New("unsupported")

// DeleteDir implements server.Driver.
func (*httpProxyDriver) DeleteDir(*server.Context, string) error {
	return errUnsupported
}

// DeleteFile implements server.Driver.
func (*httpProxyDriver) DeleteFile(*server.Context, string) error {
	return errUnsupported
}

// ListDir implements server.Driver.
func (*httpProxyDriver) ListDir(*server.Context, string, func(fs.FileInfo) error) error {
	return errUnsupported
}

// MakeDir implements server.Driver.
func (*httpProxyDriver) MakeDir(*server.Context, string) error {
	return errUnsupported
}

// PutFile implements server.Driver.
func (*httpProxyDriver) PutFile(*server.Context, string, io.Reader, int64) (int64, error) {
	return 0, errUnsupported
}

// Rename implements server.Driver.
func (*httpProxyDriver) Rename(*server.Context, string, string) error {
	return errUnsupported
}

// Stat implements server.Driver.
func (*httpProxyDriver) Stat(*server.Context, string) (fs.FileInfo, error) {
	return nil, errUnsupported
}
