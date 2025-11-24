package app

import (
	"context"
	"net/http"

	"github.com/Zaba505/infra/services/machine/endpoint"
	"github.com/Zaba505/infra/services/machine/firestore"
	"github.com/z5labs/humus/rest"
)

type Config struct {
	rest.Config `config:",squash"`
	Firestore   FirestoreConfig `config:"firestore"`
}

type FirestoreConfig struct {
	ProjectID string `config:"project_id"`
}

func Init(ctx context.Context, cfg Config) (*rest.Api, error) {
	fsClient, err := firestore.NewClient(ctx, cfg.Firestore.ProjectID)
	if err != nil {
		return nil, err
	}

	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api := rest.NewApi(
		cfg.OpenApi.Title,
		cfg.OpenApi.Version,
		rest.Liveness(healthHandler),
		rest.Readiness(healthHandler),
		endpoint.PostMachines(fsClient),
	)

	return api, nil
}
