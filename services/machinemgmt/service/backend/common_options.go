package backend

import "go.uber.org/zap"

type commonOptions struct {
	log *zap.Logger
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

func Logger(logger *zap.Logger) CommonOption {
	return commonOptionFunc(func(co *commonOptions) {
		co.log = logger
	})
}
