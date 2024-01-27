package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/pkg/framework"
	"github.com/Zaba505/infra/services/lb-sink/service"
)

//go:embed config.yaml
var cfgSrc []byte

func main() {
	framework.RunHttp(
		bytes.NewReader(cfgSrc),
		service.Init,
	)
}
