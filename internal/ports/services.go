package ports

import (
	"context"

	"drones/internal/core/domain"
)

// JWT token service
// JwtTokenService defines the interface for JWT token-related operations in the system.
// This service provides methods to generate and verify JWT tokens.
//
// Current implementation includes:
// - Token generation
// - Token verification
//
// TODO: Future enhancements should include:
// - Token revocation
// - Support for different signing algorithms
//
// For example, to handle token revocation:
//
//	 // Revoke token
//		RevokeToken(ctx context.Context, tokenString string) error
type JwtTokenService interface {
	// Generate JWT token
	GenerateToken(ctx context.Context, userID string, userType string) (string, error)

	// Verify JWT token
	VerifyToken(ctx context.Context, tokenString string) (string, string, error)
}

// Auth service
// AuthService defines the interface for authentication-related operations in the system.
// This service provides methods to handle user login and authentication processes.
//
// Current implementation includes:
// - User login handling
//
// TODO: Future enhancements should include:
// - Two-factor authentication
//
// For example, to handle user login:
//
//	 // Verify OTP
//		VerifyOtp(ctx context.Context, token string, otpCode string) (*domain.Auth, error)
type AuthService interface {

	// Handle Login
	Login(ctx context.Context, name string, userType string) (*domain.Auth, error)

	// Verify Token
	VerifyToken(ctx context.Context, tokenString string) (string, string, error)

	// Get user by ID
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}

// Users service
// UserService defines the interface for user-```go
// UserService defines the interface for user-related operations in the system.
// This service provides methods to retrieve user information by various criteria.
//
// Current implementation includes:
// - User retrieval by ID
// - User retrieval by name and type combination
//
// TODO: Future enhancements should include:
// - CreateUser: Create a new user in the system
// - UpdateUser: Update existing user information
// - DeleteUser: Remove a user from the system
// - ListUsers: Retrieve a paginated list of users
//
// For example, to create a new user:
//
//	 // Create a new user
//		CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
type UserService interface {
	// Get user by ID
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)

	// Get user by name and type
	GetUserByNameAndType(ctx context.Context, name string, userType string) (*domain.User, error)
}

// Orders service
// OrdersService defines the interface for order-related operations in the system.
// This service provides methods to manage orders including creation, retrieval, updating, deletion, and listing.
//
// Current implementation includes:
// - CreateOrder: Create a new order in the system
// - GetByID: Retrieve an order by its ID
// - UpdateOrder: Update existing order information
// - DeleteOrder: Remove an order from the system
// - ListOrders: Retrieve a paginated list of orders based on filters
//
// TODO: Future enhancements should include:
// - BulkCreateOrders: Create multiple orders in a single operation
// - CancelOrder: Cancel an existing order
//
// For example, to create a new order:
//
//	// Create a new order
type OrdersService interface {
	// Define methods for order data persistence
	CreateOrder(ctx context.Context, userID string, order *domain.CreateOrderRequest) (*domain.Order, error)

	// GetByID retrieves an order by its ID
	GetOrderByID(ctx context.Context, orderID string, options domain.OrderFilter) (*domain.Order, error)

	// GetCurrentOrder retrieves the current active order for a user
	GetOrderByFilter(ctx context.Context, options domain.OrderFilter) (*domain.Order, error)

	// Withdraw an order
	Withdraw(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)

	// Update an existing order in the repository
	UpdateOrder(ctx context.Context, orderID string, update *domain.UpdateOrderRequest, options domain.OrderFilter) (*domain.Order, error)

	// Reserve an order
	Reserve(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)

	// Confirm pickup for an order
	ConfirmPickup(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)

	// Start transit for an order
	StartTransit(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)

	// Mark order as arrived
	ConfirmArrived(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)

	// Deliverd
	ConfirmDelivery(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)

	// Delivery failed
	DeliveryFailed(ctx context.Context, orderID string, userID string, options domain.OrderFilter) (*domain.Order, error)


	// handoff an order
	Handoff(ctx context.Context, orderID string, droneID string, options domain.OrderFilter) (*domain.Order, error)

	// reassign an order
	Reassign(ctx context.Context, orderID string, droneID string, options domain.OrderFilter) (*domain.Order, error)

	// DeleteOrder deletes an order from the repository
	DeleteOrder(ctx context.Context, orderID string, options domain.OrderFilter) error

	// ListOrders retrieves a list of orders based on the provided filter
	ListOrders(ctx context.Context, options domain.PaginationOption[domain.OrderFilter]) (*domain.Pagination[domain.OrderDTO], error)

	// Order location update
	UpadateOrderLocation(ctx context.Context, userID, orderID string, currentLat, currentLon, currentAltitude float64, options domain.OrderFilter) (*domain.Order, error)
}

type DronesService interface {

	// GetDroneByID retrieves a drone by its ID
	GetDroneByID(ctx context.Context, droneID string) (*domain.Drone, error)
	// Get orders by filer
	GetDroneByFilter(ctx context.Context, options domain.DroneFilter) (*domain.Drone, error)

	// Nearby drones
	NearbyDrones(ctx context.Context, lat, lon, radiusKm float64) ([]*domain.Drone, error)

	// List drones with pagination
	ListDrones(ctx context.Context, options domain.PaginationOption[domain.DroneFilter]) (*domain.Pagination[domain.DroneDTO], error)

	// Update drone
	UpdateDrone(ctx context.Context, droneID string, update *domain.UpdateDroneRequest) (*domain.Drone, error)

	// Action broken
	UpdateDroneStatus(ctx context.Context, userID, droneID string, status domain.DroneStatus) (*domain.Drone, error)

	// Heartbeat
	ProcessHeartbeat(ctx context.Context,  droneID string , userId string, req domain.HeartbeatRequest) (*domain.Drone, error)
}
