package server

import (
	"app/pkg/utils/otelutil"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Name        string
	Listen      string
	engine      *gin.Engine
	setup       func(engine *gin.Engine)
	otelOptions otelutil.Options
}

func New(name string, listen string, opts otelutil.Options, setup func(engine *gin.Engine)) *Server {
	engine := gin.Default()
	if opts.ServiceName == "" {
		opts.ServiceName = name
	}
	return &Server{
		Name:        name,
		Listen:      listen,
		otelOptions: opts,
		engine:      engine,
		setup:       setup,
	}
}

func (server *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Set up OpenTelemetry.
	otelShutdown, err := otelutil.SetupOTelSDK(ctx, server.otelOptions)
	if err != nil {
		log.Fatalf("failed to initialize otel sdk: %v", err)
		return err
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	server.setup(server.engine)

	if InLambda() {
		ServeLambda(server.engine)
		return nil
	} else {
		srv := &http.Server{
			Addr:        server.Listen,
			BaseContext: func(_ net.Listener) context.Context { return ctx },
			Handler:     server.engine,
		}
		srvErr := make(chan error, 1)
		go func() {
			srvErr <- srv.ListenAndServe()
		}()
		// Wait for interruption.
		select {
		case err = <-srvErr:
			// Error when starting HTTP server.
			log.Fatal(err)
			return err
		case <-ctx.Done():
			// Wait for first CTRL+C.
			// Stop receiving signal notifications as soon as possible.
			stop()
		}
		err = srv.Shutdown(context.Background())
		return err
	}
}
