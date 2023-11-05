package service

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/Zaba505/infra/pkg/httpvalidate"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/z5labs/app"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type config struct {
	Http struct {
		Port uint `config:"port"`
	} `config:"http"`
}

type runtime struct {
	log *otelzap.Logger

	// http
	port uint

	started atomic.Bool
	healthy atomic.Bool
	serving atomic.Bool
}

func BuildRuntime(bc app.BuildContext) (app.Runtime, error) {
	var cfg config
	err := bc.Config.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	tp, err := configureOtel()
	if err != nil {
		logger.Error("failed to configure otel", zap.Error(err))
		return nil, err
	}
	otel.SetTracerProvider(tp)

	rt := &runtime{
		log:  otelzap.New(logger),
		port: cfg.Http.Port,
	}
	return rt, nil
}

func (rt *runtime) Run(ctx context.Context) error {
	conn, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(rt.port), 10))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	registerEndpoint(
		mux,
		"/health/startup",
		httpvalidate.Request(
			http.HandlerFunc(rt.startupHandler),
			httpvalidate.ForMethods(http.MethodGet),
		),
	)
	registerEndpoint(
		mux,
		"/health/liveness",
		httpvalidate.Request(
			http.HandlerFunc(rt.livenessHandler),
			httpvalidate.ForMethods(http.MethodGet),
		),
	)
	registerEndpoint(
		mux,
		"/health/readiness",
		httpvalidate.Request(
			http.HandlerFunc(rt.readinessHandler),
			httpvalidate.ForMethods(http.MethodGet),
		),
	)

	s := &http.Server{
		Handler: otelhttp.NewHandler(mux, "machinemgmt", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gctx.Done()
		defer func() {
			tp := otel.GetTracerProvider()
			stp, ok := tp.(interface {
				Shutdown(context.Context) error
			})
			if !ok {
				return
			}
			stp.Shutdown(context.Background())
		}()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer rt.log.Info("stopped service")

		s.Shutdown(ctx)
		return nil
	})
	g.Go(func() error {
		rt.log.Info("started serving")
		rt.started.Store(true)
		rt.healthy.Store(true)
		rt.serving.Store(true)
		return s.Serve(conn)
	})

	err = g.Wait()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func registerEndpoint(mux *http.ServeMux, path string, h http.Handler) {
	mux.Handle(
		path,
		otelhttp.WithRouteTag(path, h),
	)
}

// report whether this service is ready to begin accepting traffic
func (rt *runtime) startupHandler(w http.ResponseWriter, req *http.Request) {
	started := rt.started.Load()
	if started {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// report whether this service is healthy or needs to be restarted
func (rt *runtime) livenessHandler(w http.ResponseWriter, req *http.Request) {
	healthy := rt.healthy.Load()
	if healthy {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// report whether this service is able to accept traffic
func (rt *runtime) readinessHandler(w http.ResponseWriter, req *http.Request) {
	serving := rt.serving.Load()
	if serving {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
