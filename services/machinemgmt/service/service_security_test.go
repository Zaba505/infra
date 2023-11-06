package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

func TestRun(t *testing.T) {
	addrCh := make(chan string)
	rt := &runtime{
		log: otelzap.L(),
		listen: func(network, addr string) (net.Listener, error) {
			defer close(addrCh)

			ls, err := net.Listen(network, addr)
			if err != nil {
				return nil, err
			}
			addrCh <- ls.Addr().String()
			return ls, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		rt.Run(ctx)
	}()

	addr := <-addrCh

	t.Run("will return a 400 status code", func(t *testing.T) {
		t.Run("if the id param is missing", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://%s/bootstrap/image", addr))
			if !assert.Nil(t, err) {
				return
			}
			defer resp.Body.Close()

			if !assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
				return
			}
		})
	})

	t.Run("will return a 405 status code", func(t *testing.T) {
		endpoints := []string{
			"/health/startup",
			"/health/liveness",
			"/health/readiness",
			"/bootstrap/image",
		}

		for _, endpoint := range endpoints {
			endpoint := endpoint
			t.Run("if a non GET request is made to "+endpoint, func(t *testing.T) {
				resp, err := http.Post(fmt.Sprintf("http://%s%s", addr, endpoint), "application/json", nil)
				if !assert.Nil(t, err) {
					return
				}
				defer resp.Body.Close()

				if !assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode) {
					return
				}
			})
		}
	})
}
