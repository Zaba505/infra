package main

import (
	"bytes"
	_ "embed"

	"github.com/Zaba505/infra/services/lb-sink/service"

	"github.com/z5labs/app"
	"github.com/z5labs/app/pkg/otelconfig"
)

//go:embed config.yaml
var cfg []byte

func main() {
	app.New(
		app.Config(bytes.NewReader(cfg)),
		app.InitTracerProvider(func(bc app.BuildContext) (otelconfig.Initializer, error) {
			var cfg struct {
				OTel struct {
					GCP struct {
						ProjectId   string `config:"projectId"`
						ServiceName string `config:"serviceName"`
					} `config:"gcp"`
				} `config:"otel"`
			}
			err := bc.Config.Unmarshal(&cfg)
			if err != nil {
				return nil, err
			}

			var otelIniter otelconfig.Initializer = otelconfig.Noop
			if cfg.OTel.GCP.ProjectId != "" {
				otelIniter = otelconfig.GoogleCloud(
					otelconfig.GoogleCloudProjectId(cfg.OTel.GCP.ProjectId),
					otelconfig.ServiceName(cfg.OTel.GCP.ServiceName),
				)
			}
			return otelIniter, nil
		}),
		app.WithRuntimeBuilderFunc(service.BuildRuntime),
	).Run()
}
