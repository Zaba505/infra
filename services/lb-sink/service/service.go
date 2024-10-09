package service

import (
	"context"
	"net/http"

	"github.com/Zaba505/infra/pkg/rest"
)

type Config struct{}

func Init(ctx context.Context, cfg Config) ([]rest.Endpoint, error) {
	endpoints := []rest.Endpoint{
		rest.Get(
			"/",
			&unavailableHandler{},
			rest.StatusCode(http.StatusServiceUnavailable),
		),
	}
	return endpoints, nil
}

type unavailableHandler struct{}

func (*unavailableHandler) Handle(_ context.Context, _ *rest.Empty) (*rest.Empty, error) {
	return &rest.Empty{}, nil
}
