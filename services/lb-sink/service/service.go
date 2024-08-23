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

type Empty struct{}

func (*unavailableHandler) Handle(_ context.Context, _ *Empty) (*Empty, error) {
	// TODO: repond with 503
	return &Empty{}, nil
}
