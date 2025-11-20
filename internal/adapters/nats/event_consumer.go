package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	config "drones/configs"
	"drones/internal/core/domain"
	"drones/internal/ports"
)

// EventConsumer implements the EventConsumer interface using NATS
type EventConsumer struct {
	conn     *nats.Conn
	subs     map[string]*nats.Subscription
	handlers map[domain.EventType]ports.EventHandler
	logger   ports.Logger
	config   config.NATSConfig
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewEventConsumer creates a new NATS event consumer
func NewEventConsumer(config config.NATSConfig, logger ports.Logger) ports.EventConsumer {
	return &EventConsumer{
		subs:     make(map[string]*nats.Subscription),
		handlers: make(map[domain.EventType]ports.EventHandler),
		logger:   logger,
		config:   config,
	}
}

// Start starts consuming events
func (c *EventConsumer) Start(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	// Connect to NATS
	conn, err := nats.Connect(c.config.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	c.conn = conn

	// Subscribe to subjects we want to consume from
	consumeSubjects := map[string]string{
		"drones.events": c.config.Subjects.DronesEvents,
		"orders.events": c.config.Subjects.OrdersEvents,
		"users.events":  c.config.Subjects.UsersEvents,
	}

	for name, subject := range consumeSubjects {
		sub, err := c.conn.QueueSubscribe(subject, c.config.QueueGroup, c.messageHandler)
		if err != nil {
			return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
		}
		c.subs[name] = sub
	}

	c.logger.Info("Event consumer started",
		zap.Int("subscriptions", len(c.subs)),
		zap.String("queue_group", c.config.QueueGroup))

	return nil
}

// Stop stops consuming events
func (c *EventConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}

	c.wg.Wait()

	// Unsubscribe from all subjects
	for name, sub := range c.subs {
		if err := sub.Unsubscribe(); err != nil {
			c.logger.Error("Failed to unsubscribe",
				zap.String("subscription", name),
				zap.Error(err))
		}
	}

	// Close NATS connection
	if c.conn != nil {
		c.conn.Close()
	}

	c.logger.Info("Event consumer stopped")
	return nil
}

// RegisterHandler registers a handler for a specific event type
func (c *EventConsumer) RegisterHandler(eventType domain.EventType, handler ports.EventHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers[eventType] = handler
	c.logger.Info("Event handler registered",
		zap.String("event_type", string(eventType)))

	return nil
}

// messageHandler handles incoming NATS messages
func (c *EventConsumer) messageHandler(msg *nats.Msg) {
	select {
	case <-c.ctx.Done():
		return
	default:
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			if err := c.handleMessage(msg); err != nil {
				c.logger.Error("Failed to handle message",
					zap.String("subject", msg.Subject),
					zap.Error(err))
			}
		}()
	}
}

// handleMessage handles a NATS message
func (c *EventConsumer) handleMessage(msg *nats.Msg) error {
	// Log the received message
	c.logger.Info("Received message",
		zap.String("subject", msg.Subject),
		zap.Int("size", len(msg.Data)))

	// Parse the domain event
	var domainEvent domain.DomainEvent
	if err := json.Unmarshal(msg.Data, &domainEvent); err != nil {
		c.logger.Error("Failed to unmarshal domain event",
			zap.String("subject", msg.Subject),
			zap.Error(err))
		return fmt.Errorf("failed to unmarshal domain event: %w", err)
	} else {
		// Log the parsed domain event
		c.logger.Info("Parsed domain event",
			zap.String("event_type", string(domainEvent.Type)),
			zap.String("event_id", domainEvent.ID),
			zap.String("aggregate_id", domainEvent.AggregateID),
			zap.Time("timestamp", domainEvent.Timestamp))
	}

	// Find and execute the handler
	c.mu.RLock()
	handler, exists := c.handlers[domainEvent.Type]
	c.mu.RUnlock()

	if !exists {
		c.logger.Info("No handler found for event type",
			zap.String("event_type", string(domainEvent.Type)))
		return nil // Not an error - we just don't handle this event type
	}

	c.logger.Info("Handling event",
		zap.String("event_type", string(domainEvent.Type)),
		zap.String("event_id", domainEvent.ID),
		zap.String("aggregate_id", domainEvent.AggregateID))

	// Add correlation ID to context
	ctx := context.WithValue(c.ctx, "correlation_id", domainEvent.Metadata.CorrelationID)

	// Handle the event
	if err := handler.Handle(ctx, domainEvent); err != nil {
		return fmt.Errorf("handler failed for event %s: %w", domainEvent.Type, err)
	}

	c.logger.Info("Event handled successfully",
		zap.String("event_type", string(domainEvent.Type)),
		zap.String("event_id", domainEvent.ID))

	return nil
}
