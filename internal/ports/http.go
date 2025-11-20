package ports

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HTTPHandler defines the interface for HTTP handlers in the system.
// This interface provides methods to set up routes, handle not found cases, and perform health checks.
//
// Current implementation includes:
// - Route setup
// - Not found handling
// - Health check endpoint
//
// TODO: Future enhancements should include:
// - Static file serving
//
// For example, to serve static files:
//
//	 // Serve static files
//		ServeStaticFiles(w http.ResponseWriter, r *http.Request)
type HTTPHandler interface {
	// Setup routes
	SetupRoutes(router *mux.Router)

	// Show routes
	ShowRoutes(router *mux.Router) error

	// Not found handler
	NotFound(w http.ResponseWriter, r *http.Request)

	// Health check
	HandleHealth(w http.ResponseWriter, r *http.Request)
}

// AuthHTTPHandler defines the interface for authentication-related HTTP handlers in the system.
// This interface provides methods to handle user authorization processes.
//
// Current implementation includes:
// - User login handling
//
// TODO: Future enhancements should include:
// - Logout handling
// - User registration handling
// - Token refresh handling
//
// For example, to prevent brute-force attacks
//
//	 // User registration
//		HandleUserRegistration(w http.ResponseWriter, r *http.Request)
type AuthHTTPHandler interface {
	// User authorize
	HandleLogin(w http.ResponseWriter, r *http.Request)
}
