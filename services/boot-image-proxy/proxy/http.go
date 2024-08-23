package proxy

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
	"go.opentelemetry.io/otel"
)

type Option func(*HttpDriver)

func Logger(log *slog.Logger) Option {
	return func(hd *HttpDriver) {
		hd.log = log
	}
}

func HttpClient(c *http.Client) Option {
	return func(hd *HttpDriver) {
		hd.http = c
	}
}

func HttpTarget(s string) Option {
	return func(hd *HttpDriver) {
		hd.target = s
	}
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpDriver struct {
	log    *slog.Logger
	http   httpClient
	target string

	newRequestWithContext func(context.Context, string, string, io.Reader) (*http.Request, error)
}

func NewHttpDriver(opts ...Option) *HttpDriver {
	hd := &HttpDriver{
		http:                  http.DefaultClient,
		newRequestWithContext: http.NewRequestWithContext,
	}
	for _, opt := range opts {
		opt(hd)
	}
	return hd
}

type httpStatusError struct {
	status int
}

func (e httpStatusError) Error() string {
	return fmt.Sprintf("unexpected http status code from backend: %d", e.status)
}

// Open implements ftpserver.ClientDriver.
func (d *HttpDriver) Open(name string) (afero.File, error) {
	return d.OpenFile(name, 0, 0)
}

// OpenFile implements ftpserver.ClientDriver.
func (d *HttpDriver) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	spanCtx, span := otel.Tracer("service").Start(context.Background(), "HttpDriver.GetFile")
	defer span.End()

	endpoint := "https://" + path.Join(d.target, name)
	log := d.log.With(slog.String("endpoint", endpoint))
	log.InfoContext(spanCtx, "getting file from backend")

	req, err := d.newRequestWithContext(spanCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.ErrorContext(spanCtx, "failed to construct http request", slog.String("error", err.Error()))
		return nil, err
	}

	resp, err := d.http.Do(req)
	if err != nil {
		log.ErrorContext(spanCtx, "http request failed", slog.String("error", err.Error()))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.ErrorContext(
			spanCtx,
			"received unexpected http status code from backend",
			slog.Int("http_status_code", resp.StatusCode),
		)
		return nil, httpStatusError{status: resp.StatusCode}
	}
	defer resp.Body.Close()
	log.InfoContext(spanCtx, "received successful response from backend")

	file := mem.NewFileHandle(mem.CreateFile(name))
	n, err := io.Copy(file, resp.Body)
	if err != nil {
		log.ErrorContext(
			spanCtx,
			"failed to copy response body to inmemory file",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	file.Seek(0, 0)
	log.InfoContext(spanCtx, "read file from backend", slog.Int64("file_size_in_bytes", n))
	return file, nil
}

// Chmod implements ftpserver.ClientDriver.
func (d *HttpDriver) Chmod(name string, mode fs.FileMode) error {
	d.log.Info("chmod")
	return nil
}

// Chown implements ftpserver.ClientDriver.
func (d *HttpDriver) Chown(name string, uid int, gid int) error {
	d.log.Info("chown")
	return nil
}

// Chtimes implements ftpserver.ClientDriver.
func (d *HttpDriver) Chtimes(name string, atime time.Time, mtime time.Time) error {
	d.log.Info("chtimes")
	return nil
}

// Create implements ftpserver.ClientDriver.
func (d *HttpDriver) Create(name string) (afero.File, error) {
	d.log.Info("create")
	return nil, nil
}

// Mkdir implements ftpserver.ClientDriver.
func (d *HttpDriver) Mkdir(name string, perm fs.FileMode) error {
	d.log.Info("mkdir")
	return nil
}

// MkdirAll implements ftpserver.ClientDriver.
func (d *HttpDriver) MkdirAll(path string, perm fs.FileMode) error {
	d.log.Info("mkdir all")
	return nil
}

// Name implements ftpserver.ClientDriver.
func (d *HttpDriver) Name() string {
	d.log.Info("name")
	return ""
}

// Remove implements ftpserver.ClientDriver.
func (d *HttpDriver) Remove(name string) error {
	d.log.Info("remove")
	return nil
}

// RemoveAll implements ftpserver.ClientDriver.
func (d *HttpDriver) RemoveAll(path string) error {
	d.log.Info("remove all")
	return nil
}

// Rename implements ftpserver.ClientDriver.
func (d *HttpDriver) Rename(oldname string, newname string) error {
	d.log.Info("rename")
	return nil
}

type dirInfo struct {
	name string
}

// Stat implements ftpserver.ClientDriver.
func (d *HttpDriver) Stat(name string) (fs.FileInfo, error) {
	d.log.Info("stat", slog.String("path", name))
	return dirInfo{name: name}, nil
}

// IsDir implements fs.FileInfo.
func (dirInfo) IsDir() bool {
	return true
}

// ModTime implements fs.FileInfo.
func (dirInfo) ModTime() time.Time {
	return time.Now()
}

// Mode implements fs.FileInfo.
func (dirInfo) Mode() fs.FileMode {
	return fs.ModeDir
}

// Name implements fs.FileInfo.
func (fi dirInfo) Name() string {
	return fi.name
}

// Size implements fs.FileInfo.
func (dirInfo) Size() int64 {
	return 0
}

// Sys implements fs.FileInfo.
func (dirInfo) Sys() any {
	return nil
}
