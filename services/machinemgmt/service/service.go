package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

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
	mux.Handle(
		"/hello",
		otelhttp.WithRouteTag("/hello", http.HandlerFunc(rt.helloHandler)),
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
		return s.Serve(conn)
	})

	err = g.Wait()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (rt *runtime) helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
