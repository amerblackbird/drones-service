package ports

import (
	"context"
	"drones/internal/core/domain"
	"drones/internal/core/events"
)

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	// Start establishes connection to event system
	Start() error

	// Publish order created event
	PublishOrderCreated(ctx context.Context, event events.OrderCreatedEvent) error

	PublishOrderUpdated(ctx context.Context, event events.OrderUpdatedEvent) error

	// Drone Events
	Stop() error
}

// EventConsumer defines the interface for consuming domain events
type EventConsumer interface {
	// Start starts consuming events
	Start(ctx context.Context) error

	// Stop stops consuming events
	Stop() error

	// RegisterHandler registers a handler for a specific event type
	RegisterHandler(eventType domain.EventType, handler EventHandler) error
}

// EventHandler defines the interface for handling domain events
type EventHandler interface {
	Handle(ctx context.Context, event domain.DomainEvent) error
}
