package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Zaba505/infra/services/machinemgmt/service/backend"

	"cloud.google.com/go/storage"
	"github.com/z5labs/bedrock"
	brhttp "github.com/z5labs/bedrock/http"
	"github.com/z5labs/bedrock/http/httpvalidate"
	"github.com/z5labs/bedrock/pkg/otelslog"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	OTel struct {
		GCP struct {
			ProjectId   string `config:"projectId"`
			ServiceName string `config:"serviceName"`
		} `config:"gcp"`
	} `config:"otel"`

	Http struct {
		Port uint `config:"port"`
	} `config:"http"`

	Storage struct {
		Bucket string `config:"bucket"`
	} `config:"storage"`
}

type storageClient interface {
	GetBootstrapImage(context.Context, *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error)
}

func BuildRuntime(bc bedrock.BuildContext) (bedrock.Runtime, error) {
	var cfg config
	err := bc.Config.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			AddSource: true,
		},
	))

	gs, err := storage.NewClient(context.Background())
	if err != nil {
		logger.Error("failed to create storage client", slog.Any("error", err))
		return nil, err
	}
	bucket := gs.Bucket(cfg.Storage.Bucket)
	storageService := backend.NewStorageService(
		backend.Logger(logger.Handler()),
		backend.GoogleCloudBucket(bucket),
		backend.ObjectHasher(sha256.New),
	)

	rt := brhttp.NewRuntime(
		brhttp.ListenOnPort(cfg.Http.Port),
		brhttp.LogHandler(logger.Handler()),
		brhttp.Handle(
			"/bootstrap/image",
			httpvalidate.Request(
				http.Handler(&bootstrapImageHandler{
					log:     otelslog.New(logger.Handler()),
					storage: storageService,
				}),
				httpvalidate.ForMethods(http.MethodGet),
				httpvalidate.ExactParams("id"),
			),
		),
	)
	return rt, nil
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
