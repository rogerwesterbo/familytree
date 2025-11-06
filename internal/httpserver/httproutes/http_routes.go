package httproutes

import (
	"net/http"
	"strings"

	"github.com/rogerwesterbo/familytree/internal/httpserver/handlers/v1personshandler"
	"github.com/rogerwesterbo/familytree/internal/httpserver/handlers/v1relationshipshandler"
	"github.com/rogerwesterbo/familytree/internal/httpserver/middleware"
	_ "github.com/rogerwesterbo/familytree/internal/httpserver/swaggerdocs" // swagger docs
	"github.com/rogerwesterbo/familytree/internal/services/v1personservice"
	"github.com/rogerwesterbo/familytree/internal/services/v1ratelimitservice"
	"github.com/rogerwesterbo/familytree/internal/services/v1relationshipservice"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router holds the handlers and provides HTTP routing
type Router struct {
	mux                  *http.ServeMux
	authMiddleware       *middleware.AuthMiddleware
	corsMiddleware       *middleware.CORSMiddleware
	personsHandler       *v1personshandler.Handler
	relationshipsHandler *v1relationshipshandler.Handler
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(
	rateLimiter *v1ratelimitservice.RateLimiter,
	authMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CORSMiddleware,
	personService *v1personservice.PersonService,
	relationshipService *v1relationshipservice.RelationshipService,
) *http.ServeMux {

	// Initialize handlers with services
	personsHandler := v1personshandler.NewHandler(personService)
	relationshipsHandler := v1relationshipshandler.NewHandler(relationshipService)

	r := &Router{
		mux:                  http.NewServeMux(),
		corsMiddleware:       corsMiddleware,
		authMiddleware:       authMiddleware,
		personsHandler:       personsHandler,
		relationshipsHandler: relationshipsHandler,
	}

	r.registerRoutes()
	return r.mux
}

// registerRoutes sets up all HTTP routes
func (r *Router) registerRoutes() {
	// Swagger documentation
	r.mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// API v1 routes - wrap all API routes with a base handler that applies JSON middleware by default
	r.mux.HandleFunc("/v1/", r.v1Router)
}

// v1Router routes all /v1/* requests and applies appropriate middleware
func (r *Router) v1Router(w http.ResponseWriter, req *http.Request) {
	//path := req.URL.Path

	// Wrap handler with authentication middleware
	authenticatedHandler := r.authMiddleware.Authenticate(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		// Export endpoints use plain text middleware (exception to JSON default)

		// All other API routes use JSON middleware
		middleware.JSONContentType(r.handleAPIRoutes)(rw, request)
	}))

	authenticatedHandler.ServeHTTP(w, req)
}

// handleAPIRoutes handles the actual routing logic for API endpoints
func (r *Router) handleAPIRoutes(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Route to appropriate handler
	switch {
	case path == "/v1/persons" || strings.HasPrefix(path, "/v1/persons/"):
		r.personsHandler.HandlePersons(w, req)
	case path == "/v1/relationships" || strings.HasPrefix(path, "/v1/relationships/"):
		r.relationshipsHandler.HandleRelationships(w, req)
	default:
		http.NotFound(w, req)
	}
}
