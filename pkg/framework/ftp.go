package framework

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"net"

	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/z5labs/bedrock"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"golang.org/x/sync/errgroup"
)

//go:embed ftp_config.yaml
var ftpConfigSrc []byte

type FtpConfig struct {
	Config `config:",squash"`

	FTP struct {
		CommandPort      uint `config:"command_port"`
		PassivePortRange struct {
			Start uint `config:"start"`
			End   uint `config:"end"`
		} `config:"passive_port_range"`
	} `config:"ftp"`
}

type mainDriver struct {
	log          *slog.Logger
	cfg          FtpConfig
	clientDriver ftpserver.ClientDriver
}

func (d *mainDriver) GetSettings() (*ftpserver.Settings, error) {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", d.cfg.FTP.CommandPort))
	if err != nil {
		d.log.Error(
			"failed to listen on port",
			slogfield.Uint("port", d.cfg.FTP.CommandPort),
			slogfield.Error(err),
		)
		return nil, err
	}
	settings := &ftpserver.Settings{
		Listener:          ls,
		DisableActiveMode: true,
		PassiveTransferPortRange: &ftpserver.PortRange{
			Start: int(d.cfg.FTP.PassivePortRange.Start),
			End:   int(d.cfg.FTP.PassivePortRange.End),
		},
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

type ftpRuntime struct {
	log    *slog.Logger
	name   string
	driver ftpserver.MainDriver
}

func (rt *ftpRuntime) Run(ctx context.Context) error {
	s := ftpserver.NewFtpServer(rt.driver)

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
	rt.log.ErrorContext(ctx, "service encountered unexpected error", slogfield.Error(err))
	return err
}

func RunFTP(cfg io.Reader, f func(context.Context) (ftpserver.ClientDriver, error)) {
	run(
		bedrock.RuntimeBuilderFunc(func(ctx context.Context) (bedrock.Runtime, error) {
			log := slog.New(LogHandler())

			var cfg FtpConfig
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
				log:  log,
				name: cfg.OTel.ServiceName,
				driver: &mainDriver{
					log:          log,
					cfg:          cfg,
					clientDriver: driver,
				},
			}
			return rt, nil
		}),
		bedrock.Config(bytes.NewReader(ftpConfigSrc)),
		bedrock.Config(cfg),
	)
}
