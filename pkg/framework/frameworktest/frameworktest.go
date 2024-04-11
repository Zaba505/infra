package frameworktest

import (
	"log/slog"

	"github.com/z5labs/bedrock/pkg/noop"
)

var NoopLogger = slog.New(noop.LogHandler{})
