package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rogerwesterbo/familytree/internal/clients"
	"github.com/rogerwesterbo/familytree/internal/httpserver/httproutes"
	"github.com/rogerwesterbo/familytree/internal/httpserver/middleware"
	"github.com/rogerwesterbo/familytree/internal/services/v1ratelimitservice"
	"github.com/rogerwesterbo/familytree/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// HTTPServer represents the HTTP API server
type HTTPServer struct {
	address        string
	server         *http.Server
	rateLimiter    *v1ratelimitservice.RateLimiter
	authMiddleware *middleware.AuthMiddleware
	corsMiddleware *middleware.CORSMiddleware
	secure         bool
	tlsCertFile    string
	tlsKeyFile     string
}

// New creates a new HTTP server instance
func New(
	address string,
	rateLimiter *v1ratelimitservice.RateLimiter,
) (*HTTPServer, error) {
	// Initialize authentication middleware
	authMiddleware, err := middleware.NewAuthMiddleware()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize authentication middleware: %w", err)
	}

	// Initialize CORS middleware
	corsMiddleware := middleware.NewCORSMiddleware()

	// Get TLS configuration from environment
	secure := viper.GetBool(consts.HTTP_API_SECURE)
	tlsCertFile := viper.GetString(consts.HTTP_API_TLS_CERT_FILE)
	tlsKeyFile := viper.GetString(consts.HTTP_API_TLS_KEY_FILE)

	return &HTTPServer{
		address:        address,
		rateLimiter:    rateLimiter,
		authMiddleware: authMiddleware,
		corsMiddleware: corsMiddleware,
		secure:         secure,
		tlsCertFile:    tlsCertFile,
		tlsKeyFile:     tlsKeyFile,
	}, nil
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	// Create router with all routes
	router := httproutes.NewRouter(
		s.rateLimiter,
		s.authMiddleware,
		s.corsMiddleware,
		clients.PersonService,
		clients.RelationshipService,
	)

	// Wrap router with CORS middleware
	handler := s.corsMiddleware.Handler(router)

	// Configure HTTP server
	s.server = &http.Server{
		Addr:         s.address,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		var err error
		if s.secure {
			vlog.Infof("Starting HTTPS API server on %s (cert: %s, key: %s)", s.address, s.tlsCertFile, s.tlsKeyFile)
			err = s.server.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
		} else {
			vlog.Infof("Starting HTTP API server on %s", s.address)
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			vlog.Errorf("HTTP server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	vlog.Info("Stopping HTTP API server...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	vlog.Info("HTTP API server stopped")
	return nil
}
