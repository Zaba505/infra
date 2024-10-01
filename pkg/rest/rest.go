package rest

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/app"
	"github.com/z5labs/bedrock/pkg/config"
	"github.com/z5labs/bedrock/rest"
	"github.com/z5labs/bedrock/rest/endpoint"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Logging struct {
		Level slog.Level `config:"level"`
	} `config:"logging"`

	OTel struct {
		ServiceName    string `config:"service_name"`
		ServiceVersion string `config:"service_version"`
		GCP            struct {
			ProjectID string `config:"project_id"`
		} `config:"gcp"`
	} `config:"otel"`

	HTTP struct {
		Port uint `config:"port"`
	} `config:"http"`
}

type options struct {
	restOpts []rest.Option
}

type endpointOptions struct {
	endOpts []endpoint.Option
}

type EndpointOption func(*endpointOptions)

func StatusCode(status int) EndpointOption {
	return func(eo *endpointOptions) {
		eo.endOpts = append(eo.endOpts, endpoint.StatusCode(status))
	}
}

func Header(name string, pattern string, required bool) EndpointOption {
	return func(eo *endpointOptions) {
		eo.endOpts = append(eo.endOpts, endpoint.Headers(endpoint.Header{
			Name:     name,
			Pattern:  pattern,
			Required: required,
		}))
	}
}

func QueryParam(name string, pattern string, required bool) EndpointOption {
	return func(eo *endpointOptions) {
		eo.endOpts = append(eo.endOpts, endpoint.QueryParams(endpoint.QueryParam{
			Name:     name,
			Pattern:  pattern,
			Required: required,
		}))
	}
}

type Endpoint struct {
	method  string
	pattern string
	op      rest.Operation
}

func Get[Req, Resp any](pattern string, h endpoint.Handler[Req, Resp], opts ...EndpointOption) Endpoint {
	eo := &endpointOptions{}
	for _, opt := range opts {
		opt(eo)
	}

	return Endpoint{
		method:  http.MethodGet,
		pattern: pattern,
		op:      endpoint.NewOperation(h, eo.endOpts...),
	}
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.op.ServeHTTP(w, r)
}

//go:embed default_config.yaml
var defaultConfig []byte

func Run[T any](r io.Reader, f func(context.Context, T) ([]Endpoint, error)) {
	err := bedrock.Run[T](
		context.Background(),
		build(f),
		config.FromYaml(
			config.RenderTextTemplate(
				bytes.NewReader(defaultConfig),
				config.TemplateFunc("env", func(s string) string {
					return os.Getenv(s)
				}),
				config.TemplateFunc("default", func(v any, s string) any {
					if v == nil {
						return s
					}
					return v
				}),
			),
		),
		config.FromYaml(
			config.RenderTextTemplate(
				r,
				config.TemplateFunc("env", func(s string) string {
					return os.Getenv(s)
				}),
				config.TemplateFunc("default", func(v any, s string) any {
					if v == nil {
						return s
					}
					return v
				}),
			),
		),
	)
	if err == nil {
		return
	}
	// TODO
}

type noopHandler struct{}

type Empty struct{}

func (noopHandler) Handle(_ context.Context, _ *Empty) (*Empty, error) {
	return &Empty{}, nil
}

func build[T any](f func(context.Context, T) ([]Endpoint, error)) bedrock.AppBuilderFunc[T] {
	return func(ctx context.Context, cfg T) (bedrock.App, error) {
		endpoints, err := f(ctx, cfg)
		if err != nil {
			return nil, err
		}

		ls, err := net.Listen("tcp", ":80")
		if err != nil {
			return nil, err
		}

		opts := &options{
			restOpts: []rest.Option{
				rest.Listener(ls),
				rest.Endpoint(
					http.MethodGet,
					"/health/startup",
					endpoint.NewOperation(noopHandler{}),
				),
				rest.Endpoint(
					http.MethodGet,
					"/health/liveness",
					endpoint.NewOperation(noopHandler{}),
				),
			},
		}
		for _, endpoint := range endpoints {
			opts.restOpts = append(opts.restOpts, rest.Endpoint(endpoint.method, endpoint.pattern, endpoint.op))
		}

		var base bedrock.App = rest.NewApp(opts.restOpts...)
		base = app.WithOTel(
			base,
			app.OTelLoggerProvider(func(ctx context.Context) (log.LoggerProvider, error) {
				return nil, nil
			}),
			app.OTelTextMapPropogator(func(ctx context.Context) (propagation.TextMapPropagator, error) {
				return nil, nil
			}),
			app.OTelMeterProvider(func(ctx context.Context) (metric.MeterProvider, error) {
				return nil, nil
			}),
			app.OTelTracerProvider(func(ctx context.Context) (trace.TracerProvider, error) {
				return nil, nil
			}),
		)
		base = app.WithSignalNotifications(base, os.Interrupt, os.Kill)
		return base, nil
	}
}
