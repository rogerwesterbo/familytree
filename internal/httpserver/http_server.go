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
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// HTTPServer represents the HTTP API server
type HTTPServer struct {
	address        string
	server         *http.Server
	rateLimiter    *v1ratelimitservice.RateLimiter
	authMiddleware *middleware.AuthMiddleware
	corsMiddleware *middleware.CORSMiddleware
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

	return &HTTPServer{
		address:        address,
		rateLimiter:    rateLimiter,
		authMiddleware: authMiddleware,
		corsMiddleware: corsMiddleware,
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

	vlog.Infof("Starting HTTP API server on %s", s.address)

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
