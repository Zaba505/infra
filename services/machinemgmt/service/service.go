package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Zaba505/infra/pkg/framework"
	"github.com/Zaba505/infra/services/machinemgmt/service/backend"

	"cloud.google.com/go/storage"
	"github.com/z5labs/bedrock/http/httpvalidate"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	framework.Config `config:",squash"`

	Storage struct {
		Bucket string `config:"bucket"`
	} `config:"storage"`
}

func Init(ctx context.Context) (http.Handler, error) {
	var cfg Config
	err := framework.UnmarshalConfigFromContext(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	logHandler := framework.LogHandler()
	logger := slog.New(logHandler)

	gs, err := storage.NewClient(context.Background())
	if err != nil {
		logger.ErrorContext(ctx, "failed to create storage client", slogfield.Error(err))
		return nil, err
	}
	bucket := gs.Bucket(cfg.Storage.Bucket)
	storageService := backend.NewStorageService(
		backend.Logger(logger.Handler()),
		backend.GoogleCloudBucket(bucket),
		backend.ObjectHasher(sha256.New),
	)

	mux := http.NewServeMux()
	mux.Handle(
		"/bootstrap/image",
		httpvalidate.Request(
			http.Handler(&bootstrapImageHandler{
				log:     logger,
				storage: storageService,
			}),
			httpvalidate.ForMethods(http.MethodGet),
			httpvalidate.ExactParams("id"),
		),
	)
	return mux, nil
}

type storageClient interface {
	GetBootstrapImage(context.Context, *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error)
}

type bootstrapImageHandler struct {
	log     *slog.Logger
	storage storageClient
}

func (h *bootstrapImageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	imageId := params.Get("id")

	spanCtx, span := otel.Tracer("service").Start(req.Context(), "runtime.bootstrapImageHandler", trace.WithAttributes(
		attribute.String("image.id", imageId),
	))
	defer span.End()

	resp, err := h.storage.GetBootstrapImage(spanCtx, &backend.GetBootstrapImageRequest{
		ID: imageId,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.ErrorContext(
			spanCtx,
			"failed to get bootstrap image",
			slogfield.String("image_id", imageId),
			slogfield.Error(err),
		)
		return
	}
	defer resp.Body.Close()

	base64Hash := base64.URLEncoding.EncodeToString(resp.Hash)
	w.Header().Add("ETag", fmt.Sprintf("sha256/%s", base64Hash))
	w.Header().Add("Content-Type", "application/octet")

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.ErrorContext(
			spanCtx,
			"failed to write image to response",
			slogfield.String("image_id", imageId),
			slogfield.Error(err),
		)
		return
	}
}
