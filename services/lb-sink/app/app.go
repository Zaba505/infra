package app

import (
	"context"
	"net/http"

	"github.com/Zaba505/infra/services/lb-sink/endpoint"
	"github.com/z5labs/humus/rest"
)

type Config struct {
	rest.Config `config:",squash"`
}

func Init(ctx context.Context, cfg Config) (*rest.Api, error) {
	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api := rest.NewApi(
		cfg.OpenApi.Title,
		cfg.OpenApi.Version,
		rest.Liveness(healthHandler),
		rest.Readiness(healthHandler),
		endpoint.Unavailable(),
	)

	return api, nil
}
