package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Zaba505/infra/pkg/framework/frameworktest"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	t.Run("will return an error", func(t *testing.T) {
		t.Run("if it fails to unmarshal the config", func(t *testing.T) {
			unmarshalErr := errors.New("failed to unmarshal")
			unmarshalConfig := func(ctx context.Context, v any) error {
				return unmarshalErr
			}

			b := builder{
				unmarshalConfig: unmarshalConfig,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := b.build(ctx)
			if !assert.Equal(t, unmarshalErr, err) {
				return
			}
		})
	})

	t.Run("will return a server.Driver", func(t *testing.T) {
		t.Run("if the config is properly unmarshalled", func(t *testing.T) {
			unmarshalConfig := func(ctx context.Context, v any) error {
				return nil
			}

			b := builder{
				unmarshalConfig: unmarshalConfig,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			d, err := b.build(ctx)
			if !assert.Nil(t, err) {
				return
			}
			if !assert.NotNil(t, d) {
				return
			}
		})
	})
}

type httpClientFunc func(*http.Request) (*http.Response, error)

func (f httpClientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestHttpProxyDriver_GetFile(t *testing.T) {
	t.Run("will return an error", func(t *testing.T) {
		t.Run("if the http request fails to be created", func(t *testing.T) {
			reqErr := errors.New("failed to create request")
			newRequestWithContext := func(_ context.Context, _, _ string, _ io.Reader) (*http.Request, error) {
				return nil, reqErr
			}

			d := &httpProxyDriver{
				log:                   frameworktest.NoopLogger,
				newRequestWithContext: newRequestWithContext,
			}

			_, _, err := d.GetFile(nil, "", 0)
			if !assert.Equal(t, reqErr, err) {
				return
			}
		})

		t.Run("if the http client fails to execute the request", func(t *testing.T) {
			httpErr := errors.New("failed to execute request")
			httpClient := httpClientFunc(func(r *http.Request) (*http.Response, error) {
				return nil, httpErr
			})

			d := &httpProxyDriver{
				log:                   frameworktest.NoopLogger,
				newRequestWithContext: http.NewRequestWithContext,
				http:                  httpClient,
			}

			_, _, err := d.GetFile(nil, "", 0)
			if !assert.Equal(t, httpErr, err) {
				return
			}
		})

		t.Run("if the status code is not 200", func(t *testing.T) {
			httpClient := httpClientFunc(func(r *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
				}
				return resp, nil
			})

			d := &httpProxyDriver{
				log:                   frameworktest.NoopLogger,
				newRequestWithContext: http.NewRequestWithContext,
				http:                  httpClient,
			}

			_, _, err := d.GetFile(nil, "", 0)

			var he httpStatusError
			if !assert.ErrorAs(t, err, &he) {
				return
			}
			if !assert.NotEmpty(t, he.Error()) {
				return
			}
			if !assert.Equal(t, http.StatusInternalServerError, he.status) {
				return
			}
		})
	})

	t.Run("will return file", func(t *testing.T) {
		t.Run("if the http response succeeds", func(t *testing.T) {
			content := `hello world!`
			httpClient := httpClientFunc(func(r *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode:    http.StatusOK,
					Body:          io.NopCloser(strings.NewReader(content)),
					ContentLength: int64(len(content)),
				}
				return resp, nil
			})

			d := &httpProxyDriver{
				log:                   frameworktest.NoopLogger,
				newRequestWithContext: http.NewRequestWithContext,
				http:                  httpClient,
			}

			n, file, err := d.GetFile(nil, "", 0)
			if !assert.Nil(t, err) {
				return
			}
			if !assert.Equal(t, int64(len(content)), n) {
				return
			}
			defer file.Close()

			b, err := io.ReadAll(file)
			if !assert.Nil(t, err) {
				return
			}
			if !assert.Equal(t, content, string(b)) {
				return
			}
		})
	})
}
