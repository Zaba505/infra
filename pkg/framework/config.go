package framework

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/noop"
	"github.com/z5labs/bedrock/pkg/otelslog"
)

//go:embed base_config.yaml
var baseCfg []byte

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

	FTP struct {
		Port uint `config:"port"`
	} `config:"ftp"`
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

func UnmarshalConfigFromContext(ctx context.Context, v any) error {
	m := bedrock.ConfigFromContext(ctx)
	return m.Unmarshal(v)
}
