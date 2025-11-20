package http

import (
	"context"
	domain "drones/internal/core/domain"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// userContextKey is the key used to store and retrieve user information from the context
const userContextKey = contextKey("user")

// User represents the authenticated user information stored in the context
func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext retrieves the authenticated user information from the context
func UserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}
