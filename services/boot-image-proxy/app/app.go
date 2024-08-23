package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Zaba505/infra/pkg/ftp"
	"github.com/Zaba505/infra/services/boot-image-proxy/proxy"

	ftpserver "github.com/fclairamb/ftpserverlib"
)

type Config struct {
	ftp.Config `config:",squash"`

	Proxy struct {
		Target string `config:"target"`
	} `config:"proxy"`
}

func Init(ctx context.Context, cfg Config) (ftpserver.ClientDriver, error) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.Logging.Level,
	}))

	driver := proxy.NewHttpDriver(
		proxy.Logger(log),
		proxy.HttpClient(http.DefaultClient),
		proxy.HttpTarget(cfg.Proxy.Target),
	)
	return driver, nil
}
