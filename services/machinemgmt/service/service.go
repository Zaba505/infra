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
	"github.com/z5labs/app"
	apphttp "github.com/z5labs/app/http"
	"github.com/z5labs/app/http/httpvalidate"
	"github.com/z5labs/app/pkg/otelconfig"
	"github.com/z5labs/app/pkg/otelslog"
	"github.com/z5labs/app/pkg/slogfield"
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

func BuildRuntime(bc app.BuildContext) (app.Runtime, error) {
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
		backend.Logger(logger),
		backend.GoogleCloudBucket(bucket),
		backend.ObjectHasher(sha256.New),
	)

	var otelIniter otelconfig.Initializer = otelconfig.Noop
	if cfg.OTel.GCP.ProjectId != "" {
		otelIniter = otelconfig.GoogleCloud(
			otelconfig.ProjectId(cfg.OTel.GCP.ProjectId),
			otelconfig.ServiceName(cfg.OTel.GCP.ServiceName),
		)
	}

	rt := apphttp.NewRuntime(
		apphttp.ListenOnPort(cfg.Http.Port),
		apphttp.LogHandler(logger.Handler()),
		apphttp.TracerProvider(otelIniter),
		apphttp.Handle(
			"/bootstrap/image",
			httpvalidate.Request(
				http.Handler(&bootstrapImageHandler{
					log:     otelslog.New(logger),
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
