package app

import (
	"context"

	"github.com/z5labs/humus/rest"
)

type Config struct {
	rest.Config `config:",squash"`
}

func Init(ctx context.Context, cfg Config) (*rest.Api, error) {
	api := rest.NewApi(
		cfg.OpenApi.Title,
		cfg.OpenApi.Version,
	)

	return api, nil
}
