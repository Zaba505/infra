package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnavailableHandler(t *testing.T) {
	testCases := []struct {
		Name string
		Path string
	}{
		{Name: "root", Path: "/"},
		{Name: "sub", Path: "/hello"},
		{Name: "deep", Path: "/hello/world"},
		{Name: "params", Path: "/example?hello=world"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			h := unavailableHandler{}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com"+testCase.Path, nil)

			h.ServeHTTP(w, r)

			if !assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode) {
				return
			}
		})
	}
}
