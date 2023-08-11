package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/pxe/tftp"

	"github.com/z5labs/app"
)

//go:embed config.yaml
var cfgFile []byte

func main() {
	app.New(
		app.WithRuntimeBuilderFunc(tftp.BuildRuntime),
		app.Config(bytes.NewReader(cfgFile)),
	).Run()
}
