package ports

import (
	"context"
	"database/sql"

	"drones/internal/core/domain"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {

	// Close prepare statements
	Close() error

	// Get db
	GetDB() *sql.DB

	// Get user by ID
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)

	// Get user by name and type
	GetUserByNameAndType(ctx context.Context, name string, userType string) (*domain.User, error)
}

type OrdersRepository interface {
	// Close prepare statements
	Close() error

	// Get db
	GetDB() *sql.DB

	// Define methods for order data persistence
	CreateOrder(ctx context.Context, userID string, order *domain.CreateOrderRequest) (*domain.Order, error)

	// GetByID retrieves an order by its ID
	GetOrderByID(ctx context.Context, orderID string, options domain.OrderFilter) (*domain.Order, error)

	// GetOrderByFilter retrieves an order by filter
	GetOrderByFilter(ctx context.Context, options domain.OrderFilter) (*domain.Order, error)

	// Update an existing order in the repository
	UpdateOrder(ctx context.Context, orderID string, update *domain.UpdateOrderRequest) (*domain.Order, error)

	// DeleteOrder deletes an order from the repository
	DeleteOrder(ctx context.Context, orderID string) error

	// ListOrders retrieves a list of orders based on the provided filter
	ListOrders(ctx context.Context, options domain.PaginationOption[domain.OrderFilter]) (*domain.Pagination[domain.OrderDTO], error)

	// ReseverOrder reserves an order for a specific drone
	UpdateOrderStatus(ctx context.Context, orderID string, options domain.UpdateStatusRequest) (*domain.Order, error)
}

type DronesRepository interface {
	// Close prepare statements
	Close() error

	// Get db
	GetDB() *sql.DB

	// Create drone
	CreateDrone(ctx context.Context, drone *domain.CreateDroneRequest) (*domain.Drone, error)

	// Update drone
	UpdateDrone(ctx context.Context, droneID string, update *domain.UpdateDroneRequest) (*domain.Drone, error)

	// GetDroneByID retrieves a drone by its ID
	GetDroneByID(ctx context.Context, droneID string) (*domain.Drone, error)

	// Get by filter
	GetDroneByFilter(ctx context.Context, options domain.DroneFilter) (*domain.Drone, error)

	// NearbyDrones retrieves drones near a specific location within a given radius
	NearbyDrones(ctx context.Context, lat, lon, radiusKm float64) ([]*domain.Drone, error)

	// ListDrones retrieves a list of drones based on the provided filter
	ListDrones(ctx context.Context, options domain.PaginationOption[domain.DroneFilter]) (*domain.Pagination[domain.DroneDTO], error)

	// Update status
	UpdateStatusBroken(ctx context.Context, userID, droneID string, status domain.DroneStatus) (*domain.Drone, error)

	// Heartbeat
	ProcessHeartbeat(ctx context.Context, droneID string, userId string, req domain.HeartbeatRequest) (*domain.Drone, error)
}
