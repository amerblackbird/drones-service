package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	config "drones/configs"
	"drones/internal/core/domain"
	"drones/internal/core/events"
	"drones/internal/ports"
)

// EventPublisher implements the EventPublisher interface using NATS
type EventPublisher struct {
	conn   *nats.Conn
	logger ports.Logger
	config config.NATSConfig
}

// NewEventPublisher creates a new NATS event publisher
func NewEventPublisher(config config.NATSConfig, logger ports.Logger) ports.EventPublisher {
	return &EventPublisher{
		logger: logger,
		config: config,
	}
}

// Start connects to NATS
func (p *EventPublisher) Start() error {
	conn, err := nats.Connect(p.config.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	p.conn = conn
	return nil
}

func (p *EventPublisher) Stop() error {
	return p.Close()
}

// publishEvent publishes a domain event to NATS
func (p *EventPublisher) publishEvent(ctx context.Context, subject string, event domain.DomainEvent) error {
	if p.conn == nil {
		return fmt.Errorf("NATS connection not established")
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event",
			zap.String("event_type", string(event.Type)),
			zap.String("aggregate_id", event.AggregateID),
			zap.Error(err))
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create NATS message with headers
	msg := &nats.Msg{
		Subject: subject,
		Data:    eventJSON,
		Header: nats.Header{
			"event-type":     []string{string(event.Type)},
			"event-id":       []string{event.ID},
			"correlation-id": []string{event.Metadata.CorrelationID},
			"aggregate-id":   []string{event.AggregateID},
		},
	}

	err = p.conn.PublishMsg(msg)
	if err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("subject", subject),
			zap.String("event_type", string(event.Type)),
			zap.String("aggregate_id", event.AggregateID),
			zap.Error(err))
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("Event published successfully",
		"subject", subject,
		"event_type", string(event.Type),
		"event_id", event.ID,
		"aggregate_id", event.AggregateID)

	return nil
}

func (p *EventPublisher) PublishOrderCreated(ctx context.Context, event events.OrderCreatedEvent) error {
	domainEvent := domain.DomainEvent{
		ID:          generateEventID(),
		Type:        domain.EventTypeOrderCreated,
		AggregateID: event.OrderID,
		Version:     1,
		Data:        eventToMap(event),
		Metadata: domain.EventMetadata{
			Source:        "drones",
			CorrelationID: getCorrelationID(ctx),
		},
		Timestamp: time.Now(),
	}

	return p.publishEvent(ctx, p.config.Subjects.OrdersEvents, domainEvent)
}

func (p *EventPublisher) PublishOrderUpdated(ctx context.Context, event events.OrderUpdatedEvent) error {
	domainEvent := domain.DomainEvent{
		ID:          generateEventID(),
		Type:        domain.EventTypeOrderUpdated,
		AggregateID: event.OrderID,
		Version:     1,
		Data:        eventToMap(event),
		Metadata: domain.EventMetadata{
			Source:        "drones",
			CorrelationID: getCorrelationID(ctx),
		},
		Timestamp: time.Now(),
	}

	return p.publishEvent(ctx, p.config.Subjects.OrdersEvents, domainEvent)
}

// Close closes the NATS connection
func (p *EventPublisher) Close() error {
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}
