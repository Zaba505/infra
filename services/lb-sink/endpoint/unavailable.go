package endpoint

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
	"github.com/z5labs/humus/rest"
)

type unavailableHandler struct{}

func Unavailable() rest.ApiOption {
	return rest.Handle(
		http.MethodGet,
		rest.BasePath("/*"),
		&unavailableHandler{},
	)
}

func (h *unavailableHandler) RequestBody() openapi3.RequestBodyOrRef {
	return openapi3.RequestBodyOrRef{}
}

func (h *unavailableHandler) Responses() openapi3.Responses {
	return openapi3.Responses{
		MapOfResponseOrRefValues: map[string]openapi3.ResponseOrRef{
			"503": {
				Response: &openapi3.Response{
					Description: "Service Unavailable",
				},
			},
		},
	}
}

func (h *unavailableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}
