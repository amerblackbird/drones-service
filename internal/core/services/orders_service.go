package services

import (
	"context"
	"drones/internal/core/domain"
	"drones/internal/core/events"
	"drones/internal/ports"
	"fmt"
	"time"
)

type OrdersServiceImpl struct {
	repo           ports.OrdersRepository
	dronesService  ports.DronesService
	cacheService   ports.CacheService
	eventPublisher ports.EventPublisher
	logger         ports.Logger
}

func NewOrdersService(
	repo ports.OrdersRepository,
	dronesService ports.DronesService,
	cacheService ports.CacheService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) ports.OrdersService {
	return &OrdersServiceImpl{repo: repo, dronesService: dronesService, cacheService: cacheService, eventPublisher: eventPublisher, logger: logger}
}

func (s *OrdersServiceImpl) CreateOrder(ctx context.Context, userID string, order *domain.CreateOrderRequest) (*domain.Order, error) {
	newOrder, err := s.repo.CreateOrder(ctx, userID, order)
	if err != nil {
		return nil, err
	}

	// Cache the new order
	cacheKey := fmt.Sprintf("orders:%s:", newOrder.ID)

	if err := s.cacheService.Set(ctx, cacheKey, newOrder, 0); err != nil {
		s.logger.Error("Failed to cache new order", "key", cacheKey, "error", err)
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderCreated(ctx, events.OrderCreatedEvent{
		OrderID:            newOrder.ID,
		UserID:             userID,
		OriginAddress:      newOrder.OriginAddress,
		OriginLat:          newOrder.OriginLat,
		OriginLon:          newOrder.OriginLon,
		DestinationAddress: newOrder.DestinationAddress,
		DestinationLat:     newOrder.DestinationLat,
		DestinationLon:     newOrder.DestinationLon,
	}); err != nil {
		s.logger.Error("Failed to publish order created event", "orderID", newOrder.ID, "error", err)
	}

	return newOrder, nil
}

func (s *OrdersServiceImpl) GetOrderByID(ctx context.Context, orderID string, options domain.OrderFilter) (*domain.Order, error) {
	// Check cache first
	if options.IsEmpty() {
		var cachedOrder domain.Order
		cacheKey := fmt.Sprintf("orders:%s:", orderID)
		err := s.cacheService.Get(ctx, cacheKey, &cachedOrder)
		if err == nil {
			return &cachedOrder, nil
		}
	}

	return s.repo.GetOrderByID(ctx, orderID, options)
}

func (s *OrdersServiceImpl) ListOrders(ctx context.Context, options domain.PaginationOption[domain.OrderFilter]) (*domain.Pagination[domain.OrderDTO], error) {
	return s.repo.ListOrders(ctx, options)
}

func (s *OrdersServiceImpl) UpdateOrder(ctx context.Context, orderID string, update *domain.UpdateOrderRequest, options domain.OrderFilter) (*domain.Order, error) {
	if _, err := s.repo.GetOrderByID(ctx, orderID, options); err != nil {
		return nil, err
	}
	return s.repo.UpdateOrder(ctx, orderID, update)
}

func (s *OrdersServiceImpl) DeleteOrder(ctx context.Context, orderID string, options domain.OrderFilter) error {
	if _, err := s.repo.GetOrderByID(ctx, orderID, options); err != nil {
		return err
	}
	return s.repo.DeleteOrder(ctx, orderID)
}

func (s *OrdersServiceImpl) Withdraw(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		s.logger.Error("Failed to get order for withdrawal", "orderID", orderID, "error", err)
		return nil, err
	}

	if order.Status != domain.OrderStatusPending {
		return nil, domain.ErrWithdrawNotAllowed
	}

	if order.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrAlreadyWithdrawed
	}

	status := domain.OrderStatusCancelled
	cancelledAt := time.Now().Format(time.RFC3339)
	order, err = s.repo.UpdateOrder(ctx, orderID, &domain.UpdateOrderRequest{
		Status:      &status,
		UpdatedByID: &userID,
		CancelledAt: &cancelledAt,
	})

	if err != nil {
		s.logger.Error("Failed to withdraw order", "orderID", orderID, "error", err)
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  userID,
		Status:  order.Status,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) Reserve(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}
	if order.IsReserved() || order.DroneID != nil {
		return nil, domain.ErrAlreadyReserved
	}

	if order.Status != domain.OrderStatusPending {
		return nil, domain.ErrReserveNotAllowed
	}

	drone, err := s.dronesService.GetDroneByFilter(ctx, domain.DroneFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}

	if drone == nil {
		return nil, domain.ErrDroneNotFound
	}

	if drone.Status != domain.DroneStatusIdle {
		return nil, drone.Status.GetErr()
	}

	status := domain.OrderStatusReserved
	order, err = s.repo.UpdateOrderStatus(ctx, orderID, domain.UpdateStatusRequest{
		DroneID:     drone.ID,
		UpdatedByID: userID,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: userID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) ConfirmPickup(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		s.logger.Error("Failed to get order for confirm pickup", "orderID", orderID, "error", err)
		return nil, err
	}

	if order.Status != domain.OrderStatusReserved {
		return nil, domain.ErrConfirmNotAllowed
	}

	drone, err := s.dronesService.GetDroneByFilter(ctx, domain.DroneFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}

	status := domain.OrderStatusPickedUp
	order, err = s.repo.UpdateOrderStatus(ctx, orderID, domain.UpdateStatusRequest{
		DroneID:     drone.ID,
		UpdatedByID: userID,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		s.logger.Error("Failed to confirm pickup for order", "orderID", orderID, "error", err)
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: *order.DroneID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) StartTransit(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusPickedUp && order.Status != domain.OrderStatusReassigned {
		return nil, domain.ErrTransitNotAllowed
	}
	drone, err := s.dronesService.GetDroneByFilter(ctx, domain.DroneFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}
	status := domain.OrderStatusInTransit
	order, err = s.repo.UpdateOrderStatus(ctx, orderID, domain.UpdateStatusRequest{
		DroneID:     drone.ID,
		UpdatedByID: userID,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: drone.ID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) ConfirmArrived(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusInTransit {
		return nil, domain.ErrArriveNotAllowed
	}

	drone, err := s.dronesService.GetDroneByFilter(ctx, domain.DroneFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}
	status := domain.OrderStatusArrived
	order, err = s.repo.UpdateOrderStatus(ctx, orderID, domain.UpdateStatusRequest{
		DroneID:     drone.ID,
		UpdatedByID: userID,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: drone.ID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) ConfirmDelivery(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusArrived {
		return nil, domain.ErrDeliverNotAllowed
	}
	drone, err := s.dronesService.GetDroneByFilter(ctx, domain.DroneFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}
	status := domain.OrderStatusDelivered
	deliveredAt := time.Now().Format(time.RFC3339)
	order, err = s.repo.UpdateOrderStatus(ctx, orderID, domain.UpdateStatusRequest{
		DroneID:     drone.ID,
		UpdatedByID: userID,
		Status:      status,
		DeliveredAt: &deliveredAt,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: drone.ID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) Handoff(ctx context.Context, orderID string, droneID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusFailed {
		return nil, domain.ErrHandoffNotAllowed
	}

	status := domain.OrderStatusHandoff
	order, err = s.repo.UpdateOrder(ctx, orderID, &domain.UpdateOrderRequest{
		Status:      &status,
		UpdatedByID: &droneID,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: droneID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) Reassign(ctx context.Context, orderID string, droneID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusHandoff {
		return nil, domain.ErrReassignNotAllowed
	}

	status := domain.OrderStatusReassigned
	order, err = s.repo.UpdateOrder(ctx, orderID, &domain.UpdateOrderRequest{
		Status:             &status,
		UpdatedByID:        &droneID,
		DeliveredByDroneID: &droneID,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID: orderID,
		UserID:  order.UserID,
		Status:  order.Status,
		DroneID: droneID,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) DeliveryFailed(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusPickedUp && order.Status != domain.OrderStatusArrived && order.Status != domain.OrderStatusInTransit {
		return nil, domain.ErrDeliverFailedNotAllowed
	}
	drone, err := s.dronesService.GetDroneByFilter(ctx, domain.DroneFilter{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}
	status := domain.OrderStatusFailed
	failedAt := time.Now().Format(time.RFC3339)
	order, err = s.repo.UpdateOrderStatus(ctx, orderID, domain.UpdateStatusRequest{
		DroneID:     drone.ID,
		UpdatedByID: userID,
		Status:      status,
		FailAt:      &failedAt,
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrdersServiceImpl) UpadateOrderLocation(ctx context.Context, userID, orderID string, currentLat, currentLon, currentAltitude float64, options domain.OrderFilter) (*domain.Order, error) {
	_, err := s.repo.GetOrderByID(ctx, orderID, options)
	if err != nil {
		return nil, err
	}

	order, err := s.repo.UpdateOrder(ctx, orderID, &domain.UpdateOrderRequest{
		UpdatedByID:     &userID,
		CurrentLat:      &currentLat,
		CurrentLon:      &currentLon,
		CurrentAltitude: &currentAltitude,
	})
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := s.eventPublisher.PublishOrderUpdated(ctx, events.OrderUpdatedEvent{
		OrderID:         orderID,
		UserID:          order.UserID,
		Status:          order.Status,
		CurrentLat:      &currentLat,
		CurrentLon:      &currentLon,
		CurrentAltitude: &currentAltitude,
	}); err != nil {
		s.logger.Error("Failed to publish order withdrawn event", "orderID", orderID, "error", err)
	}

	return order, nil
}

func (s *OrdersServiceImpl) GetOrderByFilter(ctx context.Context, options domain.OrderFilter) (*domain.Order, error) {
	return s.repo.GetOrderByFilter(ctx, options)
}
