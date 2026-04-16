package main

import (
	"context"
	_ "embed"
	"os"

	"github.com/Zaba505/infra/services/machine/app"
)

func main() {
	os.Exit(app.Main(context.Background()))
}
