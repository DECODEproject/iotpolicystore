package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/DECODEproject/iotcommon/middleware"
	kitlog "github.com/go-kit/kit/log"
	twrpprom "github.com/joneskoo/twirp-serverhook-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	registry "github.com/thingful/retryable-registry-prometheus"
	policystore "github.com/thingful/twirp-policystore-go"
	ps "github.com/thingful/twirp-policystore-go"
	goji "goji.io"
	"goji.io/pat"
	"golang.org/x/crypto/acme/autocert"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/postgres"
	"github.com/DECODEproject/iotpolicystore/pkg/rpc"
	"github.com/DECODEproject/iotpolicystore/pkg/version"
)

var (
	buildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "decode",
			Subsystem: "policystore",
			Name:      "build_info",
			Help:      "Information about the current build of the service",
		}, []string{"name", "version", "build_date"},
	)
)

func init() {
	registry.MustRegister(buildInfo)
}

// Server is our custom server type.
type Server struct {
	srv     *http.Server
	logger  kitlog.Logger
	db      *postgres.DB
	store   ps.PolicyStore
	domains []string
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
	buildInfo.WithLabelValues(version.BinaryName, version.Version, version.BuildDate).Set(1)

	db := postgres.NewDB(config)

	store := rpc.NewPolicyStore(config, db)
	hooks := twrpprom.NewServerHooks(registry.DefaultRegisterer)

	twirpHandler := policystore.NewPolicyStoreServer(store, hooks)

	// create a goji multiplexer
	mux := goji.NewMux()
	mux.Handle(pat.Post(policystore.PolicyStorePathPrefix+"*"), twirpHandler)
	mux.Handle(pat.Get("/metrics"), promhttp.Handler())

	// pass mux into handlers to add mappings
	MuxHandlers(mux, db)

	// add cors middleware - note here we are enabling the default of allowing
	// requests from any domain.
	c := cors.New(cors.Options{})

	mux.Use(middleware.RequestIDMiddleware)
	mux.Use(c.Handler)

	metricsMiddleware := middleware.MetricsMiddleware("decode", "policystore", registry.DefaultRegisterer)
	mux.Use(metricsMiddleware)

	// create our http.Server instance
	srv := &http.Server{
		Addr:    config.ServerAddr,
		Handler: mux,
	}

	// return the instantiated server
	return &Server{
		srv:     srv,
		logger:  kitlog.With(config.Logger, "module", "http"),
		db:      db,
		store:   store,
		domains: config.Domains,
	}
}

// Start starts the server running. We also create a channel listening for
// interrupt signals before gracefully shutting down.
func (s *Server) Start() error {
	err := s.store.(Startable).Start()
	if err != nil {
		return err
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		s.logger.Log(
			"msg", "starting server",
			"addr", s.srv.Addr,
			"twirpPrefix", policystore.PolicyStorePathPrefix,
			"domains", strings.Join(s.domains, ","),
		)

		if s.isTLSEnabled() {
			m := &autocert.Manager{
				Cache:      s.db,
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.domains...),
			}

			s.srv.TLSConfig = m.TLSConfig()

			if err := s.srv.ListenAndServeTLS("", ""); err != nil {
				s.logger.Log("err", err)
				os.Exit(1)
			}
		} else {
			if err := s.srv.ListenAndServe(); err != nil {
				s.logger.Log("err", err)
				os.Exit(1)
			}
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

// isTLSEnabled returns true if we have a non empty domains list
func (s *Server) isTLSEnabled() bool {
	return len(s.domains) > 0
}
