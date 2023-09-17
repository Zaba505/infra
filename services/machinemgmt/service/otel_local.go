//go:build !gcp

package service

import (
	"github.com/Zaba505/infra/pkg/oteltrace"

	"go.opentelemetry.io/otel/trace"
)

func configureOtel() (trace.TracerProvider, error) {
	return oteltrace.Configure(
		oteltrace.ServiceName("machinemgmt"),
	)
}
