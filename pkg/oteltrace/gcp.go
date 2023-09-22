//go:build gcp

package oteltrace

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	projectId   string
	serviceName string
}

func (config) config() {}

func GoogleCloudProject(id string) Option {
	return func(cfg Config) {
		c := cfg.(*config)
		c.projectId = id
	}
}

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
	exporter, err := texporter.New(texporter.WithProjectID(cfg.projectId))
	if err != nil {
		return nil, err
	}

	// Identify your application using resource detection
	res, err := resource.New(
		context.Background(),
		// Use the GCP resource detector to detect information about the GCP platform
		resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create trace provider with the exporter.
	//
	// By default it uses AlwaysSample() which samples all traces.
	// In a production environment or high QPS setup please use
	// probabilistic sampling.
	// Example:
	//   tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.0001)), ...)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return tp, nil
}
