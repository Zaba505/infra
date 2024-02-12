package framework

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/z5labs/bedrock"
	bdhttp "github.com/z5labs/bedrock/http"
	"github.com/z5labs/bedrock/pkg/lifecycle"
	"github.com/z5labs/bedrock/pkg/otelconfig"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func RunHttp(cfg io.Reader, f func(context.Context) (http.Handler, error)) {
	bedrock.New(
		bedrock.Config(bytes.NewReader(baseCfg)),
		bedrock.Config(cfg),
		bedrock.Hooks(
			initLogHandler(),
			lifecycle.ManageOTel(func(ctx context.Context) (otelconfig.Initializer, error) {
				var cfg Config
				err := UnmarshalConfigFromContext(ctx, &cfg)
				if err != nil {
					return nil, err
				}

				var otelIniter otelconfig.Initializer = otelconfig.Noop
				if cfg.OTel.GCP.ProjectID != "" {
					res, err := resource.New(
						context.Background(),
						resource.WithDetectors(gcp.NewDetector()),
						resource.WithAttributes(
							semconv.ServiceName(cfg.OTel.ServiceName),
							semconv.ServiceVersion(cfg.OTel.ServiceVersion),
						),
					)
					if err != nil {
						return nil, err
					}
					res, err = resource.Merge(
						resource.Default(),
						res,
					)
					if err != nil {
						return nil, err
					}
					otelIniter = otelconfig.GoogleCloud(
						otelconfig.GoogleCloudProjectId(cfg.OTel.GCP.ProjectID),
					)
				}
				return otelIniter, nil
			}),
		),
		bedrock.WithRuntimeBuilderFunc(func(ctx context.Context) (bedrock.Runtime, error) {
			var cfg Config
			err := UnmarshalConfigFromContext(ctx, &cfg)
			if err != nil {
				return nil, err
			}

			logHandler := LogHandler()
			logger := slog.New(logHandler)

			h, err := f(ctx)
			if err != nil {
				logger.ErrorContext(ctx, "failed to initialize http handler", slogfield.Error(err))
				return nil, err
			}

			rt := bdhttp.NewRuntime(
				bdhttp.LogHandler(logHandler),
				bdhttp.ListenOnPort(cfg.HTTP.Port),
				bdhttp.Handle("/", h),
			)
			return rt, nil
		}),
	).Run()
}
