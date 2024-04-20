package framework

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/config/configtmpl"
	"github.com/z5labs/bedrock/pkg/lifecycle"
	"github.com/z5labs/bedrock/pkg/noop"
	"github.com/z5labs/bedrock/pkg/otelconfig"
	"github.com/z5labs/bedrock/pkg/otelslog"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

//go:embed framework_config.yaml
var frameworkCfgSrc []byte

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
}

func UnmarshalConfigFromContext(ctx context.Context, v any) error {
	m := bedrock.ConfigFromContext(ctx)
	return m.Unmarshal(v)
}

func run(rb bedrock.RuntimeBuilder, opts ...bedrock.Option) error {
	return bedrock.
		New(
			append(
				[]bedrock.Option{
					bedrock.Hooks(defaultHooksWith()...),
					bedrock.ConfigTemplateFunc("env", configtmpl.Env),
					bedrock.ConfigTemplateFunc("default", configtmpl.Default),
					bedrock.Config(bytes.NewReader(frameworkCfgSrc)),
					bedrock.WithRuntimeBuilder(rb),
				},
				opts...,
			)...,
		).
		Run()
}

var logHandler slog.Handler = noop.LogHandler{}

func initLogHandler() func(*bedrock.Lifecycle) {
	return func(life *bedrock.Lifecycle) {
		life.PreBuild(func(ctx context.Context) error {
			var cfg Config
			err := UnmarshalConfigFromContext(ctx, &cfg)
			if err != nil {
				fmt.Print(err)
				return err
			}

			logHandler = otelslog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     cfg.Logging.Level,
				AddSource: true,
			}))
			return nil
		})
	}
}

func LogHandler() slog.Handler {
	return logHandler
}

func initOTel() func(*bedrock.Lifecycle) {
	return lifecycle.ManageOTel(func(ctx context.Context) (otelconfig.Initializer, error) {
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
	})
}

func defaultHooksWith(fs ...func(*bedrock.Lifecycle)) []func(*bedrock.Lifecycle) {
	return append([]func(*bedrock.Lifecycle){
		initLogHandler(),
		initOTel(),
	}, fs...)
}
