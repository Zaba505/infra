package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Zaba505/infra/services/machinemgmt/service/backend"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

func TestStartupHandler(t *testing.T) {
	t.Run("will return 200 response code", func(t *testing.T) {
		t.Run("if started flag is true", func(t *testing.T) {
			rt := &runtime{}
			rt.started.Store(true)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/startup", nil)

			rt.startupHandler(w, r)

			if !assert.Equal(t, http.StatusOK, w.Result().StatusCode) {
				return
			}
		})
	})

	t.Run("will return 503 response code", func(t *testing.T) {
		t.Run("if started flag is false", func(t *testing.T) {
			rt := &runtime{}
			rt.started.Store(false)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/startup", nil)

			rt.startupHandler(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})

		t.Run("if started flag is zero", func(t *testing.T) {
			rt := &runtime{}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/startup", nil)

			rt.startupHandler(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})
	})
}

func TestLivenessHandler(t *testing.T) {
	t.Run("will return 200 response code", func(t *testing.T) {
		t.Run("if healthy flag is true", func(t *testing.T) {
			rt := &runtime{}
			rt.healthy.Store(true)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/liveness", nil)

			rt.livenessHandler(w, r)

			if !assert.Equal(t, http.StatusOK, w.Result().StatusCode) {
				return
			}
		})
	})

	t.Run("will return 503 response code", func(t *testing.T) {
		t.Run("if healthy flag is false", func(t *testing.T) {
			rt := &runtime{}
			rt.healthy.Store(false)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/liveness", nil)

			rt.livenessHandler(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})

		t.Run("if healthy flag is zero", func(t *testing.T) {
			rt := &runtime{}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/liveness", nil)

			rt.livenessHandler(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})
	})
}

func TestReadinessHandler(t *testing.T) {
	t.Run("will return 200 response code", func(t *testing.T) {
		t.Run("if healthy flag is true", func(t *testing.T) {
			rt := &runtime{}
			rt.serving.Store(true)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/readiness", nil)

			rt.readinessHandler(w, r)

			if !assert.Equal(t, http.StatusOK, w.Result().StatusCode) {
				return
			}
		})
	})

	t.Run("will return 503 response code", func(t *testing.T) {
		t.Run("if healthy flag is false", func(t *testing.T) {
			rt := &runtime{}
			rt.serving.Store(false)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/readiness", nil)

			rt.readinessHandler(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})

		t.Run("if healthy flag is zero", func(t *testing.T) {
			rt := &runtime{}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/healthy/readiness", nil)

			rt.readinessHandler(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})
	})
}

type storageClientFunc func(context.Context, *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error)

func (f storageClientFunc) GetBootstrapImage(ctx context.Context, req *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error) {
	return f(ctx, req)
}

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(b []byte) (int, error) {
	return f(b)
}

func TestBootstrapImageHandler(t *testing.T) {
	t.Run("will return a 500 status code", func(t *testing.T) {
		t.Run("if the image fails to be retrieved from storage", func(t *testing.T) {
			storageErr := errors.New("storage failed")
			storage := storageClientFunc(func(c context.Context, gbir *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error) {
				return nil, storageErr
			})

			rt := &runtime{
				log:     otelzap.L(),
				storage: storage,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com?id=1", nil)

			rt.bootstrapImageHandler(w, r)

			if !assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode) {
				return
			}
		})

		t.Run("if it fails to read from the image read closer", func(t *testing.T) {
			readErr := errors.New("read failed")
			storage := storageClientFunc(func(c context.Context, gbir *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error) {
				resp := &backend.GetBootstrapImageResponse{
					Hash: []byte("abc"),
					Body: io.NopCloser(readerFunc(func(b []byte) (int, error) {
						return 0, readErr
					})),
				}
				return resp, nil
			})

			rt := &runtime{
				log:     otelzap.L(),
				storage: storage,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com?id=1", nil)

			rt.bootstrapImageHandler(w, r)

			if !assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode) {
				return
			}
		})
	})

	t.Run("will return a 200 status code", func(t *testing.T) {
		t.Run("if the image is successfully returned from storage", func(t *testing.T) {
			storage := storageClientFunc(func(c context.Context, gbir *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error) {
				resp := &backend.GetBootstrapImageResponse{
					Hash: []byte("abc"),
					Body: io.NopCloser(strings.NewReader("hello, world")),
				}
				return resp, nil
			})

			rt := &runtime{
				log:     otelzap.L(),
				storage: storage,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com?id=1", nil)

			rt.bootstrapImageHandler(w, r)

			if !assert.Equal(t, http.StatusOK, w.Result().StatusCode) {
				return
			}
		})
	})
}
