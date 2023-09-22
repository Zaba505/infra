package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloHandler(t *testing.T) {
	t.Run("will return hello world", func(t *testing.T) {
		t.Run("if a request is sent to /hello", func(t *testing.T) {
			rt := &runtime{}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			rt.helloHandler(w, r)

			b, err := io.ReadAll(w.Body)
			if !assert.Nil(t, err) {
				return
			}

			s := string(b)
			if !assert.Equal(t, "Hello, world!", s) {
				return
			}
		})
	})
}
