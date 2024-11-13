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

	"github.com/swaggest/openapi-go/openapi3"
	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/app"
	"github.com/z5labs/bedrock/pkg/config"
	"github.com/z5labs/bedrock/rest"
	"github.com/z5labs/bedrock/rest/endpoint"
	"github.com/z5labs/bedrock/rest/mux"
	"go.opentelemetry.io/otel/log"
	lognoop "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
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

func Get[I, O any, Req endpoint.Request[I], Resp endpoint.Response[O]](pattern string, h endpoint.Handler[I, O], opts ...EndpointOption) Endpoint {
	eo := &endpointOptions{}
	for _, opt := range opts {
		opt(eo)
	}

	return Endpoint{
		method:  http.MethodGet,
		pattern: pattern,
		op:      endpoint.NewOperation[I, O, Req, Resp](h, eo.endOpts...),
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

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	log.Error("failed to run rest application", slog.String("error", err.Error()))
}

type noopHandler struct{}

type Empty struct{}

func (noopHandler) Handle(_ context.Context, _ *Empty) (*Empty, error) {
	return &Empty{}, nil
}

func (Empty) ContentType() string {
	return ""
}

func (Empty) Validate() error {
	return nil
}

func (Empty) OpenApiV3Schema() (*openapi3.Schema, error) {
	return &openapi3.Schema{}, nil
}

func (Empty) ReadRequest(r *http.Request) error {
	return nil
}

func (Empty) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
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
				rest.Register(rest.Endpoint{
					Method:    http.MethodGet,
					Pattern:   "/health/startup",
					Operation: endpoint.NewOperation(noopHandler{}),
				}),
				rest.Register(rest.Endpoint{
					Method:    http.MethodGet,
					Pattern:   "/health/liveness",
					Operation: endpoint.NewOperation(noopHandler{}),
				}),
			},
		}
		for _, endpoint := range endpoints {
			opts.restOpts = append(
				opts.restOpts,
				rest.Register(rest.Endpoint{
					Method:    mux.Method(endpoint.method),
					Pattern:   endpoint.pattern,
					Operation: endpoint.op,
				}))
		}

		var base bedrock.App = rest.NewApp(opts.restOpts...)
		base = app.WithOTel(
			base,
			app.OTelLoggerProvider(func(ctx context.Context) (log.LoggerProvider, error) {
				return lognoop.NewLoggerProvider(), nil
			}),
			app.OTelTextMapPropogator(func(ctx context.Context) (propagation.TextMapPropagator, error) {
				tmp := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
				return tmp, nil
			}),
			app.OTelMeterProvider(func(ctx context.Context) (metric.MeterProvider, error) {
				return metricnoop.NewMeterProvider(), nil
			}),
			app.OTelTracerProvider(func(ctx context.Context) (trace.TracerProvider, error) {
				return tracenoop.NewTracerProvider(), nil
			}),
		)
		base = app.WithSignalNotifications(base, os.Interrupt, os.Kill)
		return base, nil
	}
}
