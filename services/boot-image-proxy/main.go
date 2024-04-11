package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/pkg/framework"
	"github.com/Zaba505/infra/services/boot-image-proxy/service"
)

//go:embed config.yaml
var configSrc []byte

func main() {
	framework.RunFTP(bytes.NewReader(configSrc), service.Init)
}
