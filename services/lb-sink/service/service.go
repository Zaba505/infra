package service

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/z5labs/app"
	apphttp "github.com/z5labs/app/http"
	"github.com/z5labs/app/pkg/otelconfig"
)

type config struct {
	Http struct {
		Port uint `config:"port"`
	} `config:"http"`
}

func BuildRuntime(bc app.BuildContext) (app.Runtime, error) {
	var cfg config
	err := bc.Config.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	logHandler := slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			AddSource: true,
		},
	)

	rt := apphttp.NewRuntime(
		apphttp.ListenOnPort(cfg.Http.Port),
		apphttp.LogHandler(logHandler),
		apphttp.TracerProvider(otelconfig.GoogleCloud(
			otelconfig.ProjectId(os.Getenv("GOOGLE_CLOUD_PROJECT")),
			otelconfig.ServiceName("lb-sink"),
		)),
		apphttp.Handle("/", &unavailableHandler{}),
	)
	return rt, nil
}

type unavailableHandler struct{}

// report 503 Service Unavailable for all requests
func (h *unavailableHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}
