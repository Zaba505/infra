package framework

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"goftp.io/server/v2"
)

type ftpRuntime struct {
	log    *slog.Logger
	name   string
	port   uint
	driver server.Driver
}

func (rt *ftpRuntime) Run(ctx context.Context) error {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", rt.port))
	if err != nil {
		rt.log.ErrorContext(
			ctx,
			"failed to listen on port",
			slogfield.Uint("port", rt.port),
			slogfield.Error(err),
		)
		return err
	}

	s, err := server.NewServer(&server.Options{
		Driver: rt.driver,
		Logger: &server.DiscardLogger{},
		Perm:   server.NewSimplePerm("", ""),
		Name:   rt.name,
	})
	if err != nil {
		rt.log.ErrorContext(ctx, "failed to initialize ftp server", slogfield.Error(err))
		return err
	}

	err = s.Serve(ls)
	return err
}

func RunFTP(cfg io.Reader, f func(context.Context) (server.Driver, error)) {
	bedrock.
		New(
			bedrock.Config(bytes.NewReader(baseCfg)),
			bedrock.Config(cfg),
			bedrock.Hooks(
				initLogHandler(),
			),
			bedrock.WithRuntimeBuilderFunc(func(ctx context.Context) (bedrock.Runtime, error) {
				log := slog.New(LogHandler())

				var cfg Config
				err := UnmarshalConfigFromContext(ctx, &cfg)
				if err != nil {
					log.ErrorContext(ctx, "failed to unmarshal config", slogfield.Error(err))
					return nil, err
				}

				driver, err := f(ctx)
				if err != nil {
					log.ErrorContext(ctx, "failed to initialize driver", slogfield.Error(err))
					return nil, err
				}
				rt := &ftpRuntime{
					log:    log,
					name:   cfg.OTel.ServiceName,
					port:   cfg.FTP.Port,
					driver: driver,
				}
				return rt, nil
			}),
		).
		Run()
}
