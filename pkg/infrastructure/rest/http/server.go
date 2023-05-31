package http

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/tracer"
)

// Server wraps up *HTTP.Server.
type Server struct {
	config           *ServerConfig
	listener         net.Listener
	logger           *logger.Logger
	server           *http.Server
	preStartCallback func() error
}

// ServerConfig indicates how a HTTP server should be initialised.
type ServerConfig struct {
	// Name of the server
	Name string

	// Address is the TCP address to listen on.
	Address string

	// GracefulShutdownHandler is a function that runs before the HTTP server is gracefully shut down.
	GracefulShutdownHandler func() error

	// KeepAlive indicates how the HTTP server should configure the connection's keep alive.
	KeepAlive struct {
		// EnforcementPolicy is used to set keepalive enforcement policy on the server-side. Server
		// will close connection with a client that violates this policy.
		EnforcementPolicy struct {
			// MinTime is the minimum amount of time a client should wait before sending a keepalive
			// ping. By default, it is 5 * time.Second.
			MinTime time.Duration

			// If true, server allows keepalive pings even when there are no active streams(RPCs). If
			// false, and client sends ping when there are no active streams, server will send GOAWAY
			// and close the connection. By default, it is false.
			PermitWithoutStream bool
		}

		// ServerParameters is used to set keepalive and max-age parameters on the server-side.
		ServerParameters struct {
			ReadTimeout       time.Duration
			ReadHeaderTimeout time.Duration
			WriteTimeout      time.Duration
			IdleTimeout       time.Duration
			MaxHeaderBytes    int
		}
	}

	// TracerProvider is the provider that uses the exporter to push traces to the collector.
	TracerProvider *tracer.Provider

	// TracerProviderShutdownHandler is a function that shuts down the tracer's exporter/provider before
	// the HTTP server is gracefully shut down.
	TracerProviderShutdownHandler func() error

	TLSConfig *tls.Config
}

// NewServer initialises a http server.
func NewServer(c *ServerConfig, logger *logger.Logger, preStartCallback func() error, router http.Handler) (*Server, error) {
	defaultServerConfig(c)
	srv := &http.Server{
		Addr:              c.Address,
		Handler:           router,
		TLSConfig:         c.TLSConfig,
		ReadTimeout:       c.KeepAlive.ServerParameters.ReadTimeout,
		ReadHeaderTimeout: c.KeepAlive.ServerParameters.ReadHeaderTimeout,
		WriteTimeout:      c.KeepAlive.ServerParameters.WriteTimeout,
		IdleTimeout:       c.KeepAlive.ServerParameters.IdleTimeout,
		MaxHeaderBytes:    c.KeepAlive.ServerParameters.MaxHeaderBytes,
	}

	return &Server{
		c,
		nil,
		logger,
		srv,
		preStartCallback,
	}, nil
}

// Addr returns the server's network address.
func (s *Server) Addr() string {
	return s.config.Address
}

// IsHTTPS returns the flag indicating whether the server is HTTPS
func (s *Server) IsHTTPS() bool {
	return strings.HasSuffix(s.config.Address, "443")
}

// HTTPServer returns the internal HTTP server instance.
func (s *Server) HTTPServer() *http.Server {
	return s.server
}

// GracefulShutdownHandler is a function that runs before the HTTP server is gracefully shut down.
func (s *Server) GracefulShutdownHandler() error {
	return s.config.GracefulShutdownHandler()
}

// GracefulStop stops the HTTP server gracefully. It stops the server from accepting new
// connections and blocks until all the pending requests are finished.
func (s *Server) GracefulStop() error {
	// Create context with timeout to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	// Shutdown HTTP and HTTPS servers
	return s.server.Shutdown(ctx)
}

// PreStartCallback is the callback function to trigger right before the server starts running.
func (s *Server) PreStartCallback() func() error {
	return s.preStartCallback
}

// Serve accepts incoming connections on the listener lis, creating a new ServerTransport and
// service goroutine for each. The service goroutines read HTTP requests and then call the
// registered handlers to reply to them. Serve returns when lis.Accept fails with fatal errors.
// lis will be closed when this method returns. Serve will return a non-nil error unless Stop or
// GracefulStop is called.
func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.Addr())
	if err != nil {
		return err
	}
	s.listener = lis

	if s.IsHTTPS() {
		return s.server.ServeTLS(s.listener, "cert-file", "key-file")
	}
	return s.server.Serve(s.listener)
}

// Type indicates if the server is HTTP or HTTPS server.
func (s *Server) Type() string {
	if s.IsHTTPS() {
		return "HTTPS"
	}
	return "HTTP"
}

// TracerProviderShutdownHandler is a function that shuts down the tracer's exporter/provider before
// the HTTP server is gracefully shut down.
func (s *Server) TracerProviderShutdownHandler() error {
	return s.config.TracerProviderShutdownHandler()
}

func defaultServerConfig(c *ServerConfig) {
	if c.TLSConfig == nil {
		// Create a new TLS configuration with recommended parameters
		c.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		}
	}

	if c.KeepAlive.EnforcementPolicy.MinTime == 0 {
		c.KeepAlive.EnforcementPolicy.MinTime = 5 * time.Second
	}

	if c.KeepAlive.ServerParameters.IdleTimeout == 0 {
		c.KeepAlive.ServerParameters.IdleTimeout = 120 * time.Second
	}

	if c.KeepAlive.ServerParameters.ReadTimeout == 0 {
		c.KeepAlive.ServerParameters.ReadTimeout = 10 * time.Second
	}

	if c.KeepAlive.ServerParameters.WriteTimeout == 0 {
		c.KeepAlive.ServerParameters.WriteTimeout = 15 * time.Second
	}

	if c.KeepAlive.ServerParameters.ReadHeaderTimeout == 0 {
		c.KeepAlive.ServerParameters.ReadHeaderTimeout = 3 * time.Second
	}
}
