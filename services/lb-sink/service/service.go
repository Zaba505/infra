package service

import (
	"context"
	"net/http"
)

func Init(ctx context.Context) (http.Handler, error) {
	return &unavailableHandler{}, nil
}

type unavailableHandler struct{}

// report 503 Service Unavailable for all requests
func (h *unavailableHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}
