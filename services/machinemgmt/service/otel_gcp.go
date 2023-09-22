//go:build gcp

package service

import (
	"os"

	"github.com/Zaba505/infra/pkg/oteltrace"

	"go.opentelemetry.io/otel/trace"
)

func configureOtel() (trace.TracerProvider, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	return oteltrace.Configure(
		oteltrace.GoogleCloudProject(projectID),
		oteltrace.ServiceName("machinemgmt"),
	)
}
