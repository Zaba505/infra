package ftp

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"

	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/z5labs/bedrock"
	bdapp "github.com/z5labs/bedrock/pkg/app"
	"github.com/z5labs/bedrock/pkg/config"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
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

	FTP struct {
		CommandPort      uint `config:"command_port"`
		PassivePortRange struct {
			Start uint `config:"start"`
			End   uint `config:"end"`
		} `config:"passive_port_range"`
	} `config:"ftp"`
}

//go:embed default_config.yaml
var defaultConfig []byte

func Run[T any](r io.Reader, f func(context.Context, T) (ftpserver.ClientDriver, error)) {
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

type app struct {
	log    *slog.Logger
	name   string
	driver ftpserver.MainDriver
}

func build[T any](f func(context.Context, T) (ftpserver.ClientDriver, error)) bedrock.AppBuilderFunc[T] {
	return func(ctx context.Context, cfg T) (bedrock.App, error) {
		driver, err := f(ctx, cfg)
		if err != nil {
			return nil, err
		}

		var base bedrock.App = &app{
			log:  nil,
			name: "",
			driver: &mainDriver{
				cfg:          Config{}, // TODO
				clientDriver: driver,
			},
		}
		base = bdapp.WithOTel(
			base,
			bdapp.OTelLoggerProvider(func(ctx context.Context) (log.LoggerProvider, error) {
				return nil, nil
			}),
			bdapp.OTelTextMapPropogator(func(ctx context.Context) (propagation.TextMapPropagator, error) {
				return nil, nil
			}),
			bdapp.OTelMeterProvider(func(ctx context.Context) (metric.MeterProvider, error) {
				return nil, nil
			}),
			bdapp.OTelTracerProvider(func(ctx context.Context) (trace.TracerProvider, error) {
				return nil, nil
			}),
		)
		base = bdapp.WithSignalNotifications(base, os.Interrupt, os.Kill)
		return base, nil
	}
}

func (app *app) Run(ctx context.Context) error {
	s := ftpserver.NewFtpServer(app.driver)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gctx.Done()
		return s.Stop()
	})
	g.Go(func() error {
		return s.ListenAndServe()
	})

	err := g.Wait()
	if err == nil {
		return nil
	}
	app.log.ErrorContext(ctx, "service encountered unexpected error", slog.String("error", err.Error()))
	return err
}

type mainDriver struct {
	log          *slog.Logger
	cfg          Config
	clientDriver ftpserver.ClientDriver
}

func (d *mainDriver) GetSettings() (*ftpserver.Settings, error) {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", d.cfg.FTP.CommandPort))
	if err != nil {
		d.log.Error(
			"failed to listen on port",
			slog.Any("port", d.cfg.FTP.CommandPort),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	d.log.Info(
		"listening for command connections",
		slog.Any("command_port", d.cfg.FTP.CommandPort),
		slog.Group(
			"passive_port_range",
			slog.Any("start", d.cfg.FTP.PassivePortRange.Start),
			slog.Any("end", d.cfg.FTP.PassivePortRange.End),
		),
	)

	settings := &ftpserver.Settings{
		Listener:          ls,
		PublicHost:        "0.0.0.0",
		DisableActiveMode: true,
		PassiveTransferPortRange: &ftpserver.PortRange{
			Start: int(d.cfg.FTP.PassivePortRange.Start),
			End:   int(d.cfg.FTP.PassivePortRange.End),
		},
		DefaultTransferType: ftpserver.TransferTypeBinary,
	}
	return settings, nil
}

func (d *mainDriver) ClientConnected(cc ftpserver.ClientContext) (string, error) {
	d.log.Info("client connected")
	return "", nil
}

func (d *mainDriver) ClientDisconnected(cc ftpserver.ClientContext) {
	d.log.Info("client disconnected")
}

func (d *mainDriver) AuthUser(cc ftpserver.ClientContext, user, pass string) (ftpserver.ClientDriver, error) {
	return d.clientDriver, nil
}

func (d *mainDriver) GetTLSConfig() (*tls.Config, error) {
	return nil, nil
}
