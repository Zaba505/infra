package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/pkg/ftp"
	"github.com/Zaba505/infra/services/boot-image-proxy/app"
)

//go:embed config.yaml
var configSrc []byte

func main() {
	ftp.Run(bytes.NewReader(configSrc), app.Init)
}
