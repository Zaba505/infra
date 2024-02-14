package service

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/Zaba505/infra/pkg/framework"

	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
	"github.com/z5labs/bedrock/http/httpclient"
	"github.com/z5labs/bedrock/pkg/slogfield"
	"go.opentelemetry.io/otel"
)

type Config struct {
	framework.FtpConfig `config:",squash"`

	Proxy struct {
		Target string `config:"target"`
	} `config:"proxy"`
}

type builder struct {
	unmarshalConfig func(context.Context, any) error
}

func Init(ctx context.Context) (ftpserver.ClientDriver, error) {
	b := builder{
		unmarshalConfig: framework.UnmarshalConfigFromContext,
	}
	return b.build(ctx)
}

func (b builder) build(ctx context.Context) (ftpserver.ClientDriver, error) {
	logHandler := framework.LogHandler()
	log := slog.New(logHandler)

	var cfg Config
	err := b.unmarshalConfig(ctx, &cfg)
	if err != nil {
		log.ErrorContext(ctx, "failed to unmarshal config", slogfield.Error(err))
		return nil, err
	}

	d := &httpProxyDriver{
		log:                   log,
		http:                  httpclient.New(),
		target:                cfg.Proxy.Target,
		newRequestWithContext: http.NewRequestWithContext,
	}
	return d, nil
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type httpProxyDriver struct {
	log    *slog.Logger
	http   httpClient
	target string

	newRequestWithContext func(context.Context, string, string, io.Reader) (*http.Request, error)
}

type httpStatusError struct {
	status int
}

func (e httpStatusError) Error() string {
	return fmt.Sprintf("unexpected http status code from backend: %d", e.status)
}

// Open implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Open(name string) (afero.File, error) {
	return d.OpenFile(name, 0, 0)
}

// OpenFile implements ftpserver.ClientDriver.
func (d *httpProxyDriver) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	spanCtx, span := otel.Tracer("service").Start(context.Background(), "httpProxyDriver.GetFile")
	defer span.End()

	d.log.InfoContext(spanCtx, "getting file from backend", slogfield.String("path", name))

	req, err := d.newRequestWithContext(spanCtx, http.MethodGet, fmt.Sprintf("https://%s/%s", d.target, name), nil)
	if err != nil {
		d.log.ErrorContext(spanCtx, "failed to construct http request", slogfield.Error(err))
		return nil, err
	}

	resp, err := d.http.Do(req)
	if err != nil {
		d.log.ErrorContext(spanCtx, "http request failed", slogfield.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d.log.ErrorContext(
			spanCtx,
			"received unexpected http status code from backend",
			slogfield.Int("http_status_code", resp.StatusCode),
		)
		return nil, httpStatusError{status: resp.StatusCode}
	}
	defer resp.Body.Close()

	file := mem.NewFileHandle(mem.CreateFile(name))
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		d.log.ErrorContext(
			spanCtx,
			"failed to copy response body to inmemory file",
			slogfield.Error(err),
		)
		return nil, err
	}
	return file, nil
}

// Chmod implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Chmod(name string, mode fs.FileMode) error {
	d.log.Info("chmod")
	return nil
}

// Chown implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Chown(name string, uid int, gid int) error {
	d.log.Info("chown")
	return nil
}

// Chtimes implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Chtimes(name string, atime time.Time, mtime time.Time) error {
	d.log.Info("chtimes")
	return nil
}

// Create implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Create(name string) (afero.File, error) {
	d.log.Info("create")
	return nil, nil
}

// Mkdir implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Mkdir(name string, perm fs.FileMode) error {
	d.log.Info("mkdir")
	return nil
}

// MkdirAll implements ftpserver.ClientDriver.
func (d *httpProxyDriver) MkdirAll(path string, perm fs.FileMode) error {
	d.log.Info("mkdir all")
	return nil
}

// Name implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Name() string {
	d.log.Info("name")
	return ""
}

// Remove implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Remove(name string) error {
	d.log.Info("remove")
	return nil
}

// RemoveAll implements ftpserver.ClientDriver.
func (d *httpProxyDriver) RemoveAll(path string) error {
	d.log.Info("remove all")
	return nil
}

// Rename implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Rename(oldname string, newname string) error {
	d.log.Info("rename")
	return nil
}

type fileInfo struct {
	name string
}

// Stat implements ftpserver.ClientDriver.
func (d *httpProxyDriver) Stat(name string) (fs.FileInfo, error) {
	d.log.Info("stat")
	return fileInfo{name: name}, nil
}

// IsDir implements fs.FileInfo.
func (fileInfo) IsDir() bool {
	return false
}

// ModTime implements fs.FileInfo.
func (fileInfo) ModTime() time.Time {
	return time.Now()
}

// Mode implements fs.FileInfo.
func (fileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

// Name implements fs.FileInfo.
func (fi fileInfo) Name() string {
	return fi.name
}

// Size implements fs.FileInfo.
func (fileInfo) Size() int64 {
	return 10
}

// Sys implements fs.FileInfo.
func (fileInfo) Sys() any {
	return nil
}