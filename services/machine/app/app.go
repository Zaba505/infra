package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Zaba505/infra/services/machine/service"
	"github.com/z5labs/bedrock/config"
)

type Config struct {
	Firestore FirestoreConfig `config:"firestore"`
}

type FirestoreConfig struct {
	ProjectID string `config:"project_id"`
}

func ConfigFromEnv(ctx context.Context) Config {
	return Config{
		Firestore: FirestoreConfig{
			ProjectID: config.Must(ctx, config.Env("GCP_PROJECT_ID")),
		},
	}
}

func Main(ctx context.Context) int {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	sigCtx, cancel := signal.NotifyContext(ctx)
	defer cancel()

	cfg := ConfigFromEnv(sigCtx)

	fsClient, err := service.NewFirestoreClient(sigCtx, cfg.Firestore.ProjectID)
	if err != nil {
		log.ErrorContext(sigCtx, "failed to initialize firestore client", slog.Any("error", err))
		return 1
	}

	return 0
}
