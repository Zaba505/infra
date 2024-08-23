package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/pkg/rest"
	"github.com/Zaba505/infra/services/machinemgmt/app"
)

//go:embed config.yaml
var cfgSrc []byte

func main() {
	rest.Run(bytes.NewReader(cfgSrc), app.Init)
}
