package backend

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"log/slog"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/googleapis/gax-go/v2"
	"github.com/z5labs/app/pkg/otelslog"
	"github.com/z5labs/app/pkg/slogfield"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type storageOptions struct {
	*commonOptions

	bucket    *storage.BucketHandle
	newHasher func() hash.Hash
}

type StorageServiceOption interface {
	Option
	applyStorage(*storageOptions)
}

type storageServiceOptionFunc func(*storageOptions)

func (f storageServiceOptionFunc) apply(v any) {
	so := v.(*storageOptions)
	f(so)
}

func (f storageServiceOptionFunc) applyStorage(so *storageOptions) {
	f(so)
}

func (f commonOptionFunc) applyStorage(so *storageOptions) {
	f(so.commonOptions)
}

func GoogleCloudBucket(bucket *storage.BucketHandle) StorageServiceOption {
	return storageServiceOptionFunc(func(so *storageOptions) {
		so.bucket = bucket
	})
}

func ObjectHasher(newHasher func() hash.Hash) StorageServiceOption {
	return storageServiceOptionFunc(func(so *storageOptions) {
		so.newHasher = newHasher
	})
}

// StorageService
type StorageService struct {
	log     *slog.Logger
	bufPool *sync.Pool

	bucket    *storage.BucketHandle
	newHasher func() hash.Hash
}

// NewStorageService
func NewStorageService(opts ...Option) *StorageService {
	sOpts := &storageOptions{
		commonOptions: &commonOptions{
			log: otelslog.New(slog.Default()),
		},
		newHasher: sha256.New,
	}
	for _, opt := range opts {
		switch x := opt.(type) {
		case StorageServiceOption:
			x.applyStorage(sOpts)
		default:
			x.apply(sOpts)
		}
	}
	s := &StorageService{
		log: sOpts.log,
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
		s.log.ErrorContext(
			spanCtx,
			"failed to get object attributes",
			slogfield.String("image_id", req.ID),
			slogfield.Error(err),
		)
		return nil, ObjectReadError{Cause: err}
	}

	r, err := obj.NewReader(spanCtx)
	if err != nil {
		s.log.ErrorContext(
			spanCtx,
			"failed to construct object reader",
			slogfield.String("image_id", req.ID),
			slogfield.Error(err),
		)
		return nil, ObjectReadError{Cause: err}
	}
	s.log.InfoContext(
		spanCtx,
		"reading bootstrap image",
		slogfield.String("image_id", req.ID),
		slogfield.Int64("object_bytes", attrs.Size),
		slogfield.Int64("object_generation", attrs.Generation),
	)

	hasher := s.newHasher()
	crc32 := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	buf, _ := s.bufPool.Get().(*bytes.Buffer)
	n, err := copyAllAndClose(io.MultiWriter(crc32, hasher, buf), r)
	if err != nil {
		s.log.ErrorContext(
			spanCtx,
			"failed to copy object to buffer",
			slogfield.String("image_id", req.ID),
			slogfield.Int64("object_bytes", attrs.Size),
			slogfield.Error(err),
		)
		return nil, ObjectReadError{Cause: err}
	}
	if n != attrs.Size {
		s.log.ErrorContext(
			spanCtx,
			"failed to copy entire object to buffer",
			slogfield.String("image_id", req.ID),
			slogfield.Int64("object_bytes", attrs.Size),
			slogfield.Int64("copied_bytes", n),
			slogfield.Error(err),
		)
		return nil, ObjectReadError{Cause: io.ErrShortWrite}
	}
	checksum := crc32.Sum32()
	if checksum != attrs.CRC32C {
		s.log.ErrorContext(
			spanCtx,
			"crc32 checksum mismatch",
			slogfield.String("image_id", req.ID),
			slogfield.Uint32("computed_checksum", checksum),
			slogfield.Uint32("object_checksum", attrs.CRC32C),
			slogfield.Error(err),
		)
		return nil, ChecksumMismatchError{}
	}
	s.log.InfoContext(
		spanCtx,
		"read bootstrap image",
		slogfield.String("image_id", req.ID),
		slogfield.Int64("object_bytes", attrs.Size),
		slogfield.Int64("object_generation", attrs.Generation),
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
