package bootstrap

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/Zaba505/infra/pkg/rest"
	"github.com/Zaba505/infra/services/machinemgmt/backend"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/z5labs/bedrock/rest/endpoint"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Option func(*handler)

func Logger(log *slog.Logger) Option {
	return func(h *handler) {
		h.log = log
	}
}

func StorageService(s *backend.StorageService) Option {
	return func(h *handler) {
		h.storage = s
	}
}

type storageClient interface {
	GetBootstrapImage(context.Context, *backend.GetBootstrapImageRequest) (*backend.GetBootstrapImageResponse, error)
}

type handler struct {
	log     *slog.Logger
	storage storageClient
}

func Endpoint(opts ...Option) rest.Endpoint {
	h := &handler{}
	for _, opt := range opts {
		opt(h)
	}
	return rest.Get(
		"/bootstrap/image/{id}",
		h,
		rest.PathParam("id", "", false),
	)
}

type Response struct {
	src io.ReadCloser
}

func (Response) ContentType() string {
	return "application/octet"
}

func (Response) OpenApiV3Schema() (*openapi3.Schema, error) {
	return &openapi3.Schema{}, nil
}

func (resp *Response) WriteTo(w io.Writer) (int64, error) {
	defer resp.src.Close()
	return io.Copy(w, resp.src)
}

func (h *handler) Handle(ctx context.Context, _ *rest.Empty) (*Response, error) {
	spanCtx, span := otel.Tracer("service").Start(ctx, "runtime.bootstrapImageHandler")
	defer span.End()

	imageId := endpoint.PathValue(ctx, "id")
	if imageId == "." || imageId == "/" {
		// TODO: Set custom response status code
		return nil, errors.New("TODO")
	}
	span.SetAttributes(attribute.String("image.id", imageId))

	resp, err := h.storage.GetBootstrapImage(spanCtx, &backend.GetBootstrapImageRequest{
		ID: imageId,
	})
	if err != nil {
		h.log.ErrorContext(
			spanCtx,
			"failed to get bootstrap image",
			slog.String("image_id", imageId),
			slog.String("error", err.Error()),
		)

		// TODO
		return nil, err
	}

	base64Hash := base64.URLEncoding.EncodeToString(resp.Hash)
	endpoint.SetResponseHeader(ctx, "ETag", fmt.Sprintf("sha256/%s", base64Hash))
	endpoint.SetResponseHeader(ctx, "Content-Type", "application/octet")

	return &Response{src: resp.Body}, nil
}
