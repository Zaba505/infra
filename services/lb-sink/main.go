package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/services/lb-sink/app"

	"github.com/z5labs/humus/rest"
)

//go:embed config.yaml
var configBytes []byte

func main() {
	rest.Run(bytes.NewReader(configBytes), app.Init)
}
