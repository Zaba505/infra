package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
