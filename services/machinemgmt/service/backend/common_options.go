package backend

import "go.uber.org/zap"

type commonOptions struct {
	log *zap.Logger
}

type Option interface {
	apply(any)
}

type CommonOption func(*commonOptions)

func (f CommonOption) apply(v any) {
	co := v.(*commonOptions)
	f(co)
}

func Logger(logger *zap.Logger) CommonOption {
	return func(co *commonOptions) {
		co.log = logger
	}
}
