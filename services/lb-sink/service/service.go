package service

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/z5labs/bedrock"
	brhttp "github.com/z5labs/bedrock/http"
)

type config struct {
	OTel struct {
		GCP struct {
			ProjectId   string `config:"projectId"`
			ServiceName string `config:"serviceName"`
		} `config:"gcp"`
	} `config:"otel"`

	Http struct {
		Port uint `config:"port"`
	} `config:"http"`
}

func BuildRuntime(bc bedrock.BuildContext) (bedrock.Runtime, error) {
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

	rt := brhttp.NewRuntime(
		brhttp.ListenOnPort(cfg.Http.Port),
		brhttp.LogHandler(logHandler),
		brhttp.Handle("/", &unavailableHandler{}),
	)
	return rt, nil
}

type unavailableHandler struct{}

// report 503 Service Unavailable for all requests
func (h *unavailableHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}
