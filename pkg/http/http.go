package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	kitlog "github.com/go-kit/kit/log"
	ps "github.com/thingful/twirp-policystore-go"
	goji "goji.io"

	"github.com/thingful/iotpolicystore/pkg/config"
	"github.com/thingful/iotpolicystore/pkg/middleware"
	"github.com/thingful/iotpolicystore/pkg/rpc"
)

// Server is our custom server type.
type Server struct {
	srv    *http.Server
	logger kitlog.Logger
	store  ps.PolicyStore
}

// Startable is an interface for a component that can be started.
type Startable interface {
	Start() error
}

// Stoppable is an interface for a component that can be stopped.
type Stoppable interface {
	Stop() error
}

// NewServer returns a new simple HTTP server.
func NewServer(config *config.Config) *Server {
	store := rpc.NewPolicyStore(config)

	// create a goji multiplexer
	mux := goji.NewMux()

	// pass mux into handlers to add mappings
	MuxHandlers(mux)

	mux.Use(middleware.RequestIDMiddleware)

	// create our http.Server instance
	srv := &http.Server{
		Addr:    config.ServerAddr,
		Handler: mux,
	}

	// return the instantiated server
	return &Server{
		srv:    srv,
		logger: kitlog.With(config.Logger, "module", "http"),
		store:  store,
	}
}

// Start starts the server running. We also create a channel listening for
// interrupt signals before gracefully shutting down.
func (s *Server) Start() error {
	err := s.store.(Startable).Start()
	if err != nil {
		return err
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		s.logger.Log(
			"msg", "starting server",
			"addr", s.srv.Addr,
		)

		if err := s.srv.ListenAndServe(); err != nil {
			s.logger.Log("err", err)
			os.Exit(1)
		}
	}()

	<-stopChan
	s.logger.Log("msg", "server is stopping")

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	return s.srv.Shutdown(ctx)
}

// Stop stops any child components, and then cleanly stops the server running
func (s *Server) Stop() error {
	s.logger.Log("msg", "stopping server")
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	err := s.store.(Stoppable).Stop()
	if err != nil {
		return err
	}

	return s.srv.Shutdown(ctx)
}
