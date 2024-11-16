package rest

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/app"
	"github.com/z5labs/bedrock/pkg/appbuilder"
	"github.com/z5labs/bedrock/pkg/config"
	"github.com/z5labs/bedrock/rest"
	"github.com/z5labs/bedrock/rest/endpoint"
	"github.com/z5labs/bedrock/rest/mux"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

type Config struct {
	Logging struct {
		Level slog.Level `config:"level"`
	} `config:"logging"`

	OTel OTelConfig `config:"otel"`

	HTTP HttpConfig `config:"http"`
}

type OTelConfig struct {
	ServiceName    string `config:"service_name"`
	ServiceVersion string `config:"service_version"`

	Trace struct {
		Enabled      bool          `config:"enabled"`
		Sampling     float64       `config:"sampling"`
		BatchTimeout time.Duration `config:"batch_timeout"`
	} `config:"trace"`

	Metric struct {
		Enabled      bool          `config:"enabled"`
		ExportPeriod time.Duration `config:"export_period"`
	} `config:"metric"`

	Log struct {
		BatchTimeout time.Duration `config:"batch_timeout"`
	} `config:"log"`
}

type HttpConfig struct {
	Port uint `config:"port"`
}

func Logger(name string) *slog.Logger {
	return otelslog.NewLogger(name)
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

func PathParam(name string, pattern string, required bool) EndpointOption {
	return func(eo *endpointOptions) {
		eo.endOpts = append(eo.endOpts, endpoint.PathParams(endpoint.PathParam{
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

func configSource(r io.Reader) config.Source {
	return config.FromYaml(
		config.RenderTextTemplate(
			r,
			config.TemplateFunc("env", func(key string) any {
				v, ok := os.LookupEnv(key)
				if ok {
					return v
				}
				return nil
			}),
			config.TemplateFunc("default", func(def, v any) any {
				if v == nil {
					return def
				}
				return v
			}),
		),
	)
}

func Run[T any](r io.Reader, f func(context.Context, T) ([]Endpoint, error)) {
	srcs := []config.Source{
		configSource(bytes.NewReader(defaultConfig)),
		configSource(r),
	}

	runner := runner{
		srcs:           srcs,
		detectResource: detectResource,
		newTraceExporter: func(ctx context.Context, oc OTelConfig) (sdktrace.SpanExporter, error) {
			return texporter.New()
		},
		newMetricExporter: func(ctx context.Context, oc OTelConfig) (sdkmetric.Exporter, error) {
			return mexporter.New()
		},
		newLogExporter: func(ctx context.Context, oc OTelConfig) (sdklog.Exporter, error) {
			return stdoutlog.New(stdoutlog.WithWriter(os.Stdout))
		},
	}

	err := run(runner, f)
	if err == nil {
		return
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	log.Error("failed to run rest application", slog.String("error", err.Error()))
}

type runner struct {
	srcs    []config.Source
	postRun postRun

	detectResource    func(context.Context, OTelConfig) (*resource.Resource, error)
	newTraceExporter  func(context.Context, OTelConfig) (sdktrace.SpanExporter, error)
	newMetricExporter func(context.Context, OTelConfig) (sdkmetric.Exporter, error)
	newLogExporter    func(context.Context, OTelConfig) (sdklog.Exporter, error)
}

func run[T any](r runner, build func(context.Context, T) ([]Endpoint, error)) error {
	m, err := config.Read(r.srcs...)
	if err != nil {
		return err
	}

	var cfg Config
	err = m.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return bedrock.Run(
		context.Background(),
		appbuilder.WithOTel(
			appbuilder.Recover(
				buildApp(build, cfg.HTTP, &r.postRun),
			),
			appbuilder.OTelTextMapPropogator(initTextMapPropogator(cfg.OTel)),
			appbuilder.OTelTracerProvider(r.initTracerProvider(cfg.OTel, &r.postRun)),
			appbuilder.OTelMeterProvider(r.initMeterProvider(cfg.OTel, &r.postRun)),
			appbuilder.OTelLoggerProvider(r.initLoggerProvider(cfg.OTel, &r.postRun)),
		),
		m,
	)
}

type postRun struct {
	hooks []app.LifecycleHook
}

func buildApp[T any](f func(context.Context, T) ([]Endpoint, error), httpCfg HttpConfig, postRun *postRun) bedrock.AppBuilder[T] {
	return bedrock.AppBuilderFunc[T](func(ctx context.Context, cfg T) (bedrock.App, error) {
		spanCtx, span := otel.Tracer("rest").Start(ctx, "buildApp.Build")
		defer span.End()

		endpoints, err := f(spanCtx, cfg)
		if err != nil {
			return nil, err
		}

		restOpts := []rest.Option{
			rest.Listener(nil), // TODO
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
		}
		for _, endpoint := range endpoints {
			restOpts = append(restOpts, rest.Register(rest.Endpoint{
				Method:    mux.Method(endpoint.method),
				Pattern:   endpoint.pattern,
				Operation: endpoint.op,
			}))
		}

		ls, err := net.Listen("tcp", fmt.Sprintf(":%d", httpCfg.Port))
		if err != nil {
			return nil, err
		}
		restOpts = append(restOpts, rest.Listener(ls))

		return app.WithSignalNotifications(
			app.WithLifecycleHooks(
				app.Recover(
					rest.NewApp(restOpts...),
				),
				app.Lifecycle{
					PostRun: composeLifecycleHooks(postRun.hooks...),
				},
			),
			os.Interrupt,
			os.Kill,
		), nil
	})
}

func composeLifecycleHooks(hooks ...app.LifecycleHook) app.LifecycleHook {
	return app.LifecycleHookFunc(func(ctx context.Context) error {
		var hookErrs []error
		for _, hook := range hooks {
			err := hook.Run(ctx)
			if err != nil {
				hookErrs = append(hookErrs, err)
			}
		}
		if len(hookErrs) == 0 {
			return nil
		}
		return errors.Join(hookErrs...)
	})
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

func initTextMapPropogator(_ OTelConfig) func(context.Context) (propagation.TextMapPropagator, error) {
	return func(ctx context.Context) (propagation.TextMapPropagator, error) {
		tmp := propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		)
		return tmp, nil
	}
}

func (r runner) initTracerProvider(cfg OTelConfig, postRun *postRun) func(context.Context) (trace.TracerProvider, error) {
	return func(ctx context.Context) (trace.TracerProvider, error) {
		if !cfg.Trace.Enabled {
			return tracenoop.NewTracerProvider(), nil
		}

		rsc, err := r.detectResource(ctx, cfg)
		if err != nil {
			return nil, err
		}

		exp, err := r.newTraceExporter(ctx, cfg)
		if err != nil {
			return nil, err
		}

		sampler := sdktrace.TraceIDRatioBased(cfg.Trace.Sampling)

		bsp := sdktrace.NewBatchSpanProcessor(
			exp,
			sdktrace.WithBatchTimeout(cfg.Trace.BatchTimeout),
		)

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithResource(rsc),
			sdktrace.WithSampler(sampler),
			sdktrace.WithSpanProcessor(bsp),
		)
		postRun.hooks = append(postRun.hooks, shutdownHook(tp))
		return tp, nil
	}
}

func (r runner) initMeterProvider(cfg OTelConfig, postRun *postRun) func(context.Context) (metric.MeterProvider, error) {
	return func(ctx context.Context) (metric.MeterProvider, error) {
		if !cfg.Metric.Enabled {
			return metricnoop.NewMeterProvider(), nil
		}

		rsc, err := r.detectResource(ctx, cfg)
		if err != nil {
			return nil, err
		}

		exp, err := r.newMetricExporter(ctx, cfg)
		if err != nil {
			return nil, err
		}

		reader := sdkmetric.NewPeriodicReader(
			exp,
			sdkmetric.WithInterval(cfg.Metric.ExportPeriod),
		)

		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(rsc),
			sdkmetric.WithReader(reader),
		)
		postRun.hooks = append(postRun.hooks, shutdownHook(mp))
		return mp, nil
	}
}

func (r runner) initLoggerProvider(cfg OTelConfig, postRun *postRun) func(context.Context) (log.LoggerProvider, error) {
	return func(ctx context.Context) (log.LoggerProvider, error) {
		rsc, err := r.detectResource(ctx, cfg)
		if err != nil {
			return nil, err
		}

		exp, err := r.newLogExporter(ctx, cfg)
		if err != nil {
			return nil, err
		}

		p := sdklog.NewBatchProcessor(
			exp,
			sdklog.WithExportInterval(cfg.Log.BatchTimeout),
		)

		lp := sdklog.NewLoggerProvider(
			sdklog.WithResource(rsc),
			sdklog.WithProcessor(p),
		)
		postRun.hooks = append(postRun.hooks, shutdownHook(lp))
		return lp, nil
	}
}

type resourceDetectorFunc func(context.Context) (*resource.Resource, error)

func (f resourceDetectorFunc) Detect(ctx context.Context) (*resource.Resource, error) {
	return f(ctx)
}

func detectResource(ctx context.Context, cfg OTelConfig) (*resource.Resource, error) {
	return resource.Detect(
		ctx,
		resourceDetectorFunc(func(ctx context.Context) (*resource.Resource, error) {
			return resource.Default(), nil
		}),
		resource.StringDetector(semconv.SchemaURL, semconv.ServiceNameKey, func() (string, error) {
			return cfg.ServiceName, nil
		}),
		resource.StringDetector(semconv.SchemaURL, semconv.ServiceVersionKey, func() (string, error) {
			return cfg.ServiceVersion, nil
		}),
		gcp.NewDetector(),
	)
}

type shutdowner interface {
	Shutdown(context.Context) error
}

func shutdownHook(s shutdowner) app.LifecycleHook {
	return app.LifecycleHookFunc(func(ctx context.Context) error {
		return s.Shutdown(ctx)
	})
}
