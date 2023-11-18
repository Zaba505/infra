package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/services/lb-sink/service"
	"github.com/z5labs/app"
)

//go:embed config.yaml
var cfg []byte

func main() {
	app.New(
		app.Config(bytes.NewReader(cfg)),
		app.WithRuntimeBuilderFunc(service.BuildRuntime),
	).Run()
}
