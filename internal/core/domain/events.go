package domain

import (
	"time"
)

type EventType string

const (
	EventTypeDroneLocationUpdated EventType = "drone_location_updated"

	// Order Events
	EventTypeOrderCreated EventType = "order_created"
	EventTypeOrderUpdated EventType = "order_updated"
)

// DomainEvent represents a domain event
type DomainEvent struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	AggregateID string                 `json:"aggregate_id"`
	Version     int                    `json:"version"`
	Data        map[string]interface{} `json:"data"`
	Metadata    EventMetadata          `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// EventMetadata contains metadata about the event
type EventMetadata struct {
	Source        string `json:"source"`
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"`
	UserID        string `json:"user_id,omitempty"`
}
