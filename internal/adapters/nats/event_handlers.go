package nats

import (
	"context"
	"encoding/json"

	"drones/internal/core/domain"
	"drones/internal/core/events"
	"drones/internal/ports"
)

// EventHandlers manages all event handlers
type EventHandlers struct {
	// Logger for logging events
	logger ports.Logger

	// Individual handlers
	authEventsHandler *AuthEventsEventHandler
}

// NewEventHandlers creates a new event handlers manager
func NewEventHandlers(dronesService ports.DronesService, logger ports.Logger) *EventHandlers {
	return &EventHandlers{
		logger:            logger,
		authEventsHandler: NewAuthEventHandler(dronesService, logger),
	}
}

// RegisterHandlers registers all event handlers with the consumer
func (h *EventHandlers) RegisterHandlers(consumer ports.EventConsumer) error {
	// Register user suspended event handler
	if err := consumer.RegisterHandler(domain.EventTypeDroneLocationUpdated, h.authEventsHandler); err != nil {
		return err
	}

	h.logger.Info("All event handlers registered successfully")
	return nil
}

type AuthEventsEventHandler struct {
	dronesService ports.DronesService
	logger        ports.Logger
}

// NewAuthEventHandler creates a new auth event handler
func NewAuthEventHandler(dronesService ports.DronesService, logger ports.Logger) *AuthEventsEventHandler {
	return &AuthEventsEventHandler{
		dronesService: dronesService,
		logger:        logger,
	}
}

// Handle handles auth events
func (h *AuthEventsEventHandler) Handle(ctx context.Context, event domain.DomainEvent) error {
	h.logger.Info("Handling auth event",
		"event_type", string(event.Type),
		"event_id", event.ID,
		"aggregate_id", event.AggregateID)

	switch event.Type {
	case domain.EventTypeDroneLocationUpdated:
		return h.handleDroneLocationUpdated(ctx, event)
	default:
		h.logger.Debug("Unhandled auth event type", "event_type", string(event.Type))
		return nil
	}
}

func (h *AuthEventsEventHandler) handleDroneLocationUpdated(ctx context.Context, event domain.DomainEvent) error {
	h.logger.Info("Processing DroneLocationUpdated event",
		"event_id", event.ID,
		"drone_id", event.AggregateID)

	// Here you would add the logic to handle the drone location update
	// For example, updating the drone's location in the database
	var request events.HeartbeatRequest
	data, err := json.Marshal(event.Data)
	if err != nil {
		h.logger.Error("Failed to marshal event data",
			"event_type", string(event.Type),
			"event_id", event.ID,
			"aggregate_id", event.AggregateID,
			"error", err)
		return err
	}
	if err := json.Unmarshal(data, &request); err != nil {
		h.logger.Error("Failed to unmarshal event data",
			"event_type", string(event.Type),
			"event_id", event.ID,
			"aggregate_id", event.AggregateID,
			"error", err)
		return err
	}

	if _, err := h.dronesService.ProcessHeartbeat(ctx, request.DroneID, request.UserID, domain.HeartbeatRequest{
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Altitude:  request.Altitude,
		Battery:   request.Battery,
	}); err != nil {
		h.logger.Error("Failed to process drone heartbeat",
			"event_type", string(event.Type),
			"event_id", event.ID,
			"aggregate_id", event.AggregateID,
			"error", err)
		return err
	}

	h.logger.Info("DroneLocationUpdated event processed successfully",
		"event_id", event.ID,
		"drone_id", event.AggregateID)

	return nil
}
