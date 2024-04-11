package framework

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net/http"

	"github.com/z5labs/bedrock"
	bdhttp "github.com/z5labs/bedrock/http"
	"github.com/z5labs/bedrock/pkg/slogfield"
)

//go:embed http_config.yaml
var httpConfigSrc []byte

type HttpConfig struct {
	Config `config:",squash"`

	HTTP struct {
		Port uint `config:"port"`
	} `config:"http"`
}

func RunHttp(cfg io.Reader, f func(context.Context) (http.Handler, error)) {
	run(
		bedrock.RuntimeBuilderFunc(func(ctx context.Context) (bedrock.Runtime, error) {
			var cfg HttpConfig
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
		bedrock.Config(bytes.NewReader(httpConfigSrc)),
		bedrock.Config(cfg),
	)
}
