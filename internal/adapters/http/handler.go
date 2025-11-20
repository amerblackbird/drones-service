package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	domain "drones/internal/core/domain"
	ports "drones/internal/ports"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// HTTPHandler implements the HTTP adapter
type HTTPHandler struct {
	authService    ports.AuthService
	ordersService  ports.OrdersService
	dronesService  ports.DronesService
	eventPublisher ports.EventPublisher
	logger         ports.Logger
	Validator      *validator.Validate
	apiPrefix      string
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	authService ports.AuthService,
	ordersService ports.OrdersService,
	dronesService ports.DronesService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
	apiPrefix string,
) ports.HTTPHandler {
	return &HTTPHandler{
		authService:    authService,
		ordersService:  ordersService,
		dronesService:  dronesService,
		eventPublisher: eventPublisher,
		logger:         logger,
		Validator:      domain.NewValidator(),
		apiPrefix:      apiPrefix,
	}
}
func (h *HTTPHandler) SetupRoutes(r *mux.Router) {
	// Middlewares
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use((&LoggingMiddleware{logger: h.logger}).log)
	r.Use(CorsMiddleware)

	r.HandleFunc("/health", h.HandleHealth).Methods("GET")

	// Auth routes
	authHandler := NewAuthHandler(h.authService, h.Validator, h.logger)
	authRouter := r.PathPrefix(fmt.Sprintf("%s/authorize", h.apiPrefix)).Subrouter()
	authHandler.RegisterRoutes(authRouter)

	// Orders routes
	ordersHandler := NewOrdersHandler(h.ordersService, h.eventPublisher, h.logger)
	ordersRouter := r.PathPrefix(fmt.Sprintf("%s/orders", h.apiPrefix)).Subrouter()
	ordersRouter.Use(func(next http.Handler) http.Handler {
		return AuthenticateMiddleware(next, "*", h.authService)
	})
	ordersHandler.RegisterRoutes(ordersRouter)

	dronesHandler := NewDronesHandler(h.dronesService, h.eventPublisher, h.logger)
	dronesRouter := r.PathPrefix(fmt.Sprintf("%s/drones", h.apiPrefix)).Subrouter()
	dronesRouter.Use(func(next http.Handler) http.Handler {
		return AuthenticateMiddleware(next, "*", h.authService)
	})
	dronesHandler.RegisterRoutes(dronesRouter)

	// TODO: Implement audit and activity logs handlers
	// auditLogsHandler := NewAuditLogsHandler(h.logger)
	// auditLogsRouter := r.PathPrefix(h.apiPrefix + "/audit-logs").Subrouter()
	// auditLogsHandler.RegisterRoutes(auditLogsRouter)

	// activityLogsHandler := NewActivityLogsHandler(h.eventPublisher, h.logger)
	// activityLogsRouter := r.PathPrefix(h.apiPrefix + "/activity-logs").Subrouter()
	// activityLogsHandler.RegisterRoutes(activityLogsRouter)

	// Not found handler
	r.NotFoundHandler = http.HandlerFunc(h.NotFound)

	// Log all routes
	if err := h.ShowRoutes(r); err != nil {
		h.logger.Error("Failed to show routes", zap.Error(err))
	}
}

func (h *HTTPHandler) ShowRoutes(r *mux.Router) error {

	var routes []string
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			// If no methods are defined, skip this route or use a default
			routes = append(routes, "* "+pathTemplate)
		} else {
			routes = append(routes, strings.Join(methods, ",")+" "+pathTemplate)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, route := range routes {
		h.logger.Info(route)
	}
	return nil
}

func (h *HTTPHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	meta := domain.ExtractRequestInfo(r)
	h.logger.Info("Meestttttttt")
	h.logger.Info("Health check", zap.Any("data", meta))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","version":"1.0.0"}`))
}

func (h *HTTPHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	meta := domain.ExtractRequestInfo(r)
	h.logger.Info("Not found", zap.Any("data", meta))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Resource not found", "service":"drones-service","version":"1.0.0"}`))
}
