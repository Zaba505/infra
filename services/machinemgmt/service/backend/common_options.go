package backend

import (
	"log/slog"

	"github.com/z5labs/bedrock/pkg/otelslog"
)

type commonOptions struct {
	log *slog.Logger
}

type Option interface {
	apply(any)
}

type CommonOption interface {
	Option
	applyCommon(*commonOptions)
}

type commonOptionFunc func(*commonOptions)

func (f commonOptionFunc) apply(v any) {
	co := v.(*commonOptions)
	f(co)
}

func (f commonOptionFunc) applyCommon(co *commonOptions) {
	f(co)
}

func Logger(h slog.Handler) CommonOption {
	return commonOptionFunc(func(co *commonOptions) {
		co.log = otelslog.New(h)
	})
}
