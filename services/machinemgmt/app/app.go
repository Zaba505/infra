package app

import (
	"context"
	"crypto/sha256"
	"log/slog"

	"github.com/Zaba505/infra/pkg/rest"
	"github.com/Zaba505/infra/services/machinemgmt/backend"
	"github.com/Zaba505/infra/services/machinemgmt/bootstrap"

	"cloud.google.com/go/storage"
)

type Config struct {
	rest.Config `config:",squash"`

	Storage struct {
		Bucket string `config:"bucket"`
	} `config:"storage"`
}

func Init(ctx context.Context, cfg Config) ([]rest.Endpoint, error) {
	log := rest.Logger("machinemgmt")

	gs, err := storage.NewClient(context.Background())
	if err != nil {
		log.ErrorContext(ctx, "failed to create storage client", slog.String("error", err.Error()))
		return nil, err
	}

	bucket := gs.Bucket(cfg.Storage.Bucket)
	storageService := backend.NewStorageService(
		backend.Logger(log.Handler()),
		backend.GoogleCloudBucket(bucket),
		backend.ObjectHasher(sha256.New),
	)

	endpoints := []rest.Endpoint{
		bootstrap.Endpoint(
			bootstrap.Logger(log),
			bootstrap.StorageService(storageService),
		),
	}
	return endpoints, nil
}
