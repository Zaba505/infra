package service

import (
	"context"
	"errors"
	"testing"
	"time"

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
