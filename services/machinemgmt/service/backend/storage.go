package backend

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/googleapis/gax-go/v2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type storageOptions struct {
	*commonOptions

	bucket    *storage.BucketHandle
	newHasher func() hash.Hash
}

type StorageServiceOption func(*storageOptions)

func (f StorageServiceOption) apply(v any) {
	so := v.(*storageOptions)
	f(so)
}

func GoogleCloudBucket(bucket *storage.BucketHandle) StorageServiceOption {
	return func(so *storageOptions) {
		so.bucket = bucket
	}
}

func ObjectHasher(newHasher func() hash.Hash) StorageServiceOption {
	return func(so *storageOptions) {
		so.newHasher = newHasher
	}
}

type objectHandle interface {
	NewReader(context.Context) (*storage.Reader, error)
}

// StorageService
type StorageService struct {
	log     *otelzap.Logger
	bufPool *sync.Pool

	bucket    *storage.BucketHandle
	newHasher func() hash.Hash
}

// NewStorageService
func NewStorageService(opts ...Option) *StorageService {
	sOpts := &storageOptions{
		commonOptions: &commonOptions{
			log: zap.NewNop(),
		},
		newHasher: sha256.New,
	}
	for _, opt := range opts {
		opt.apply(sOpts)
	}
	s := &StorageService{
		log: otelzap.New(sOpts.log),
		bufPool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
		bucket:    sOpts.bucket,
		newHasher: sOpts.newHasher,
	}

	return s
}

type GetBootstrapImageRequest struct {
	ID string
}

type GetBootstrapImageResponse struct {
	Body io.ReadCloser
	Hash []byte
}

func (s *StorageService) GetBootstrapImage(ctx context.Context, req *GetBootstrapImageRequest) (*GetBootstrapImageResponse, error) {
	spanCtx, span := otel.Tracer("backend").Start(ctx, "StorageService.GetBootstrapImage", trace.WithAttributes(
		attribute.String("image.id", req.ID),
	))
	defer span.End()

	// TODO: should prolly sanitize req.Version somehow
	obj := s.bucket.Object(fmt.Sprintf("bootstrap/%s", req.ID))
	obj.Retryer(
		storage.WithBackoff(gax.Backoff{ // TODO: parameterize this config
			Initial: 2 * time.Second,
		}),
		storage.WithPolicy(storage.RetryIdempotent),
	)

	attrs, err := obj.Attrs(spanCtx)
	if err != nil { // TODO: check if error tells object doesn't exist
		s.log.Ctx(spanCtx).Error("failed to get object attributes", zap.String("image_id", req.ID), zap.Error(err))
		return nil, ObjectReadError{Cause: err}
	}

	r, err := obj.NewReader(spanCtx)
	if err != nil {
		s.log.Ctx(spanCtx).Error("failed to construct object reader", zap.String("image_id", req.ID), zap.Error(err))
		return nil, ObjectReadError{Cause: err}
	}
	s.log.Ctx(spanCtx).Info(
		"reading bootstrap image",
		zap.String("image_id", req.ID),
		zap.Int64("object_bytes", attrs.Size),
		zap.Int64("object_generation", attrs.Generation),
	)

	hasher := s.newHasher()
	crc32 := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	buf, _ := s.bufPool.Get().(*bytes.Buffer)
	n, err := copyAllAndClose(io.MultiWriter(crc32, hasher, buf), r)
	if err != nil {
		s.log.Ctx(spanCtx).Error(
			"failed to copy object to buffer",
			zap.String("image_id", req.ID),
			zap.Int64("object_bytes", attrs.Size),
			zap.Error(err),
		)
		return nil, ObjectReadError{Cause: err}
	}
	if n != attrs.Size {
		s.log.Ctx(spanCtx).Error(
			"failed to copy entire object to buffer",
			zap.String("image_id", req.ID),
			zap.Int64("object_bytes", attrs.Size),
			zap.Int64("copied_bytes", n),
			zap.Error(err),
		)
		return nil, ObjectReadError{Cause: io.ErrShortWrite}
	}
	checksum := crc32.Sum32()
	if checksum != attrs.CRC32C {
		s.log.Ctx(spanCtx).Error(
			"crc32 checksum mismatch",
			zap.String("image_id", req.ID),
			zap.Uint32("computed_checksum", checksum),
			zap.Uint32("object_checksum", attrs.CRC32C),
			zap.Error(err),
		)
		return nil, ChecksumMismatchError{}
	}
	s.log.Ctx(spanCtx).Info(
		"read bootstrap image",
		zap.String("image_id", req.ID),
		zap.Int64("object_bytes", attrs.Size),
		zap.Int64("object_generation", attrs.Generation),
	)

	resp := &GetBootstrapImageResponse{
		Hash: hasher.Sum(nil),
		Body: &bufPoolReadCloser{
			buf:     buf,
			bufPool: s.bufPool,
		},
	}
	return resp, nil
}

type bufPoolReadCloser struct {
	buf     *bytes.Buffer
	bufPool *sync.Pool
}

func (r *bufPoolReadCloser) Read(b []byte) (int, error) {
	return r.buf.Read(b)
}

func (r *bufPoolReadCloser) Close() error {
	r.buf.Reset()
	r.bufPool.Put(r.buf)
	r.buf = nil
	return nil
}

func copyAllAndClose(dst io.Writer, src io.ReadCloser) (int64, error) {
	n, err := io.Copy(dst, src)
	if err == nil {
		return n, src.Close()
	}
	defer src.Close()
	return n, err
}
