package http

import (
	"context"
	"net/http"
	"strings"

	domain "drones/internal/core/domain"
	"drones/internal/ports"
)

type LoggingMiddleware struct {
	logger ports.Logger
}

// loggingMiddleware logs incoming requests
func (l *LoggingMiddleware) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.logger.Info("Incoming request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (h *HTTPHandler) AuthMiddleware(next http.Handler, accessRole string) http.Handler {
	var wildCard string = "*"
	if accessRole == "" {
		wildCard = accessRole
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("HandlerFunc", "HandlerFunc")
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
				Code:    domain.InvalidAuthTokenFormatError,
				Message: "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
				Code:    domain.InvalidAuthTokenFormatError,
				Message: "Invalid authorization token format",
			})
			return
		}

		tokenString := parts[1]

		// Validate the token
		userID, userType, err := h.authService.VerifyToken(r.Context(), tokenString)
		if err != nil {
			h.logger.Error("Err", err)
			h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
				Code:    domain.AuthTokenInvalidError,
				Message: "Invalid token",
			})
			return
		}

		// TODO: Add user service to handler to validate user
		// user, err := h.usersService.GetUserByID(r.Context(), userID)
		// For now, just validate token type
		if wildCard != "*" && userType != wildCard {
			h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
				Code:    domain.InvalidAuthTokenTypeError,
				Message: "Invalid token type",
			})
			return
		}

		// Add user ID to the request context for now
		ctx := context.WithValue(r.Context(), "userID", userID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticateMiddleware(next http.Handler, accessRole string, authService ports.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.InvalidAuthTokenFormatError,
				Message: "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.InvalidAuthTokenFormatError,
				Message: "Invalid authorization token format",
			})
			return
		}

		tokenString := parts[1]

		// Validate the token
		userID, userType, err := authService.VerifyToken(r.Context(), tokenString)
		if err != nil {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.AuthTokenInvalidError,
				Message: "Invalid token",
			})
			return
		}

		user, err := authService.GetUserByID(r.Context(), userID)
		if err != nil {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.UserNotFoundError,
				Message: "User not found",
			})
			return
		}

		// TODO: Add user service to handler to validate user
		// user, err := h.usersService.GetUserByID(r.Context(), userID)
		// For now, just validate token type
		if accessRole != "*" && userType != accessRole {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.InvalidAuthTokenTypeError,
				Message: "Unauthorized access for this user type",
			})
			return
		}

		// Add user ID to the request context for now
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = WithUser(r.Context(), user)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func EndUserGuard(next http.Handler) http.Handler {
	// Read user from context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || user == nil {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.UnauthenticatedError,
				Message: "User not authenticated",
			})
			return
		}
		if user.Type != "enduser" {
			ResponseWithCustomError(w, http.StatusForbidden, domain.DomainError{
				Code:    domain.AccessDeniedError,
				Message: "Access denied for this user type",
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func AdminGuard(next http.Handler) http.Handler {
	// Read user from context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || user == nil {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.UnauthenticatedError,
				Message: "User not authenticated",
			})
			return
		}
		if user.Type != "admin" {
			ResponseWithCustomError(w, http.StatusForbidden, domain.DomainError{
				Code:    domain.AccessDeniedError,
				Message: "Access denied for this user type",
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func DroneGuard(next http.Handler) http.Handler {
	// Read user from context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || user == nil {
			ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
				Code:    domain.UnauthenticatedError,
				Message: "User not authenticated",
			})
			return
		}
		if user.Type != "drone" {
			ResponseWithCustomError(w, http.StatusForbidden, domain.DomainError{
				Code:    domain.AccessDeniedError,
				Message: "Access denied for this user type",
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}


func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
