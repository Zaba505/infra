package app

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Zaba505/infra/services/machine/endpoint"
	"github.com/Zaba505/infra/services/machine/service"
	"github.com/go-chi/chi/v5"
	"github.com/sourcegraph/conc/pool"
	"github.com/z5labs/bedrock/config"
)

type Config struct {
	HTTP      HTTPConfig
	Firestore FirestoreConfig
}

type HTTPConfig struct {
	Port int
}

type FirestoreConfig struct {
	ProjectID string
}

func ConfigFromEnv(ctx context.Context) Config {
	return Config{
		HTTP: HTTPConfig{
			Port: config.Must(
				ctx,
				config.Default(
					8080,
					config.IntFromString(config.Env("HTTP_PORT")),
				),
			),
		},
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

	mux := chi.NewRouter()
	endpoint.RegisterMachines(mux, fsClient)

	srv := &http.Server{
		Handler: mux,
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.HTTP.Port))
	if err != nil {
		log.ErrorContext(sigCtx, "failed to listen on port", slog.Int("port", cfg.HTTP.Port), slog.Any("error", err))
		return 1
	}

	cert, err := generateSelfSignedCert()
	if err != nil {
		log.ErrorContext(sigCtx, "failed to generate tls certificate", slog.Any("error", err))
		return 1
	}

	ls = tls.NewListener(ls, &tls.Config{Certificates: []tls.Certificate{cert}})

	pool := pool.New().WithErrors().WithContext(sigCtx)
	pool.Go(func(ctx context.Context) error {
		log.InfoContext(ctx, "starting HTTP server", slog.Int("port", cfg.HTTP.Port))
		if err := srv.Serve(ls); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})
	pool.Go(func(ctx context.Context) error {
		<-ctx.Done()
		log.InfoContext(ctx, "shutting down HTTP server")
		return srv.Shutdown(context.Background())
	})

	if err := pool.Wait(); err != nil {
		log.ErrorContext(sigCtx, "server error", slog.Any("error", err))
		return 1
	}

	return 0
}

func generateSelfSignedCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "machine-service"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		DNSNames:     []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return tls.X509KeyPair(certPEM, keyPEM)
}
