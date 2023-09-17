//go:build !gcp

package oteltrace

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	serviceName string
}

func (config) config() {}

func ServiceName(name string) Option {
	return func(cfg Config) {
		c := cfg.(*config)
		c.serviceName = name
	}
}

func Configure(opts ...Option) (trace.TracerProvider, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	exporter, err := stdouttrace.New()
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.serviceName),
		)),
	)
	return tp, nil
}
