# Drone Delivery Management System

A comprehensive backend service for managing drone deliveries, built with Go using Hexagonal Architecture (Ports and Adapters) combined with Domain-Driven Design (DDD) principles. The system handles order management, drone fleet coordination, real-time tracking, and role-based access control.

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters) combined with **Domain-Driven Design (DDD)** principles, ensuring a clean separation between business logic and external concerns.

### Architectural Principles

**Hexagonal Architecture (Ports and Adapters)**

- **Core Domain**: Business logic is isolated in the center, independent of external frameworks and infrastructure
- **Ports**: Interfaces that define how the application core communicates with the outside world
- **Adapters**: Concrete implementations of ports that handle external integrations (HTTP, Database, Messaging, etc.)
- **Dependency Rule**: Dependencies point inward - external layers depend on inner layers, never the reverse

**Domain-Driven Design (DDD)**

- **Domain Layer**: Contains business entities, value objects, and domain logic
- **Domain Events**: Capture important business occurrences (order created, drone status changed)
- **Services**: Orchestrate domain operations and coordinate between aggregates
- **Repository Pattern**: Abstract data persistence concerns
- **Ubiquitous Language**: Code reflects business terminology (Order, Drone, Delivery, etc.)

### Layer Structure

```txt
drones/
├── cmd/app/              # Application entry point (composition root)
├── configs/              # Configuration management
├── internal/
│   ├── core/             # 🔵 DOMAIN CORE (Inner Hexagon)
│   │   ├── domain/       # Domain entities, value objects, aggregates
│   │   │   ├── order.go          # Order aggregate root
│   │   │   ├── drone.go          # Drone aggregate root
│   │   │   ├── user.go           # User entity
│   │   │   ├── events.go         # Domain events
│   │   │   └── errors.go         # Domain-specific errors
│   │   ├── services/     # Application/Domain services
│   │   │   ├── orders_service.go     # Order business logic
│   │   │   ├── drones_service.go     # Drone business logic
│   │   │   └── auth_service.go       # Authentication logic
│   │   └── events/       # Domain event definitions
│   │
│   ├── ports/            # 🔌 PORTS (Interfaces)
│   │   ├── services.go       # Inbound ports (use cases)
│   │   ├── repositories.go   # Outbound ports (data)
│   │   ├── events.go         # Outbound ports (messaging)
│   │   ├── cache.go          # Outbound ports (caching)
│   │   └── logger.go         # Outbound ports (logging)
│   │
│   └── adapters/         # 🔌 ADAPTERS (Implementations)
│       ├── http/         # Primary/Driving adapter (REST API)
│       │   ├── handler.go
│       │   ├── auth_handler.go
│       │   ├── orders_handler.go
│       │   ├── drones_handler.go
│       │   └── middlewares.go
│       ├── postgres/     # Secondary/Driven adapter (persistence)
│       │   ├── orders_repository.go
│       │   ├── drones_repository.go
│       │   └── users_repository.go
│       ├── redis/        # Secondary/Driven adapter (caching)
│       ├── nats/         # Secondary/Driven adapter (messaging)
│       │   ├── event_publisher.go
│       │   ├── event_consumer.go
│       │   └── event_handlers.go
│       └── logger/       # Secondary/Driven adapter (logging)
│
├── pkg/utils/            # Shared utilities (infrastructure concerns)
├── migrations/           # Database migrations
└── docs/                 # Documentation
```

### Key DDD Concepts

**Aggregates**

- **Order Aggregate**: Order entity with status, location, drone assignment
- **Drone Aggregate**: Drone entity with status, location, battery, assignments
- **User Aggregate**: User entity with authentication and role information

**Value Objects**

- Location (latitude, longitude, altitude)
- OrderStatus, DroneStatus
- Pagination, Filters

**Domain Events**

- `OrderCreated`: Published when a new order is submitted
- `OrderStatusChanged`: Published when order status transitions
- `DroneLocationUpdated`: Published on drone heartbeat
- `DroneStatusChanged`: Published when drone status changes

**Repository Pattern**

- Abstracts data persistence behind `ports.Repository` interfaces
- Implemented by `adapters/postgres` for PostgreSQL
- Allows switching databases without touching domain logic

**Domain Services**

- Orchestrate complex business operations across aggregates
- Enforce business rules and invariants
- Coordinate between repositories and event publishers

## Features

### Authentication & Authorization

- JWT-based authentication with self-signed tokens
- Role-based access control (Admin, Enduser, Drone)
- Bearer token authentication
- User type-specific endpoints

### Order Management

- Create, update, and track delivery orders
- Real-time order status updates
- Origin and destination management
- Order withdrawal for unpicked orders
- Bulk order retrieval for admins
- ETA and location tracking

### Drone Fleet Management

- Drone registration and identification
- Real-time location updates (lat/lon/altitude)
- Battery and payload capacity tracking
- Status management (idle, loading, delivering, returning, charging, broken, maintenance)
- Automatic order handoff on drone failure
- Maintenance scheduling

### Order Status Workflow

```
pending → reserved → picked_up → in_transit → arrived → delivered
                                      ↓    ↑
                                handoff → reassigned
```

### Drone Status Workflow

```
idle → loading → delivering → returning → idle
  ↓       ↓          ↓          ↓
charging/broken (from any state)
            ↓
under_repair → maintenanced → idle
```

## 🛠️ Technology Stack

- **Language**: Go 1.23.1
- **Web Framework**: Gorilla Mux
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Messaging**: NATS (with JetStream)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Validation**: go-playground/validator/v10
- **Logging**: Uber Zap
- **Migrations**: golang-migrate
- **Containerization**: Docker & Docker Compose

## Database Schema

### Core Tables

- **users**: User accounts with roles (admin, enduser, drone)
- **drones**: Drone fleet with specifications and status
- **orders**: Delivery orders with origin/destination
- **audit_logs**: System-wide audit trail
- **activity_logs**: User activity tracking

### Key Relationships

- `orders.user_id` → `users.id` (order owner)
- `orders.drone_id` → `drones.id` (assigned drone)
- `orders.delivered_by_drone_id` → `drones.id` (delivery completion)
- `drones.user_id` → `users.id` (drone operator)

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.23.1+ (for local development)

### Quick Start

1. **Clone the repository**

```bash
git clone <repository-url>
cd drones
```

2. **Start all services**

```bash
docker-compose up -d
```

This will start:

- PostgreSQL (port 5432)
- Redis (port 6379)
- NATS (port 4222)
- API Server (port 8080)
- PgAdmin (port 5050)

3. **Verify services**

```bash
docker-compose ps
```

4. **Access the API**

```bash
curl http://localhost:8080/health
```

### Database Management

**Run migrations manually:**

```bash
docker-compose up postgres-migrate
```

**Access PgAdmin:**

- URL: <http://localhost:5050>
- Email: <admin@admin.com>
- Type: admin

### Local Development

1. **Install dependencies**

```bash
go mod download
```

2. **Run tests**

```bash
go test ./... -v
```

##  API Documentation

### Authentication

**Login**

```http
POST /auth/login
Content-Type: application/json

{
  "name": "user_name",
  "type": "drone|enduser|admin"
}

Response: {
  "token": "jwt_token_here"
}
```

All authenticated endpoints require:

```http
Authorization: Bearer <jwt_token>
```

### Drone Endpoints

**Reserve Order**

```http
POST /drones/orders/{orderId}/reserve
```

**Pickup Order**

```http
POST /drones/orders/{orderId}/pickup
```

**Update Location (Heartbeat)**

```http
POST /drones/location
{
  "lat": 24.7136,
  "lng": 46.6753,
  "alti": 100.0,
  "speed_kmh": 45.5
}
```

**Get Assigned Order**

```http
GET /drones/orders/assigned
```

**Mark as Broken**

```http
POST /drones/broken
```

### Enduser Endpoints

**Create Order**

```http
POST /orders
{
  "origin_address": "123 Main St",
  "origin_lat": 24.7136,
  "origin_lon": 46.6753,
  "destination_address": "456 Oak Ave",
  "destination_lat": 24.7256,
  "destination_lon": 46.6853,
  "package_weight_kg": 2.5,
  "receiver_name": "John Doe",
  "receiver_phone": "+966501234567"
}
```

**Withdraw Order**

```http
POST /orders/{orderId}/withdraw
```

**Get Order Details**

```http
GET /orders/{orderId}
```

**List My Orders**

```http
GET /orders?user_id={userId}
```

### Admin Endpoints

**List All Orders**

```http
GET /admin/orders?page=1&limit=20
```

**Update Order Origin/Destination**

```http
PUT /admin/orders/{orderId}
{
  "origin_address": "New Address",
  "destination_lat": 24.7256,
  "destination_lon": 46.6853
}
```

**List Drones**

```http
GET /admin/drones?status=idle&page=1&limit=20
```

**Mark Drone as Broken/Fixed**

```http
PUT /admin/drones/{droneId}/status
{
  "status": "broken|idle"
}
```

## Testing

The project includes comprehensive test coverage:

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./pkg/utils/ -v
go test ./internal/core/services/ -v

# Run benchmarks
go test ./pkg/utils/ -bench=. -benchmem
```

### Test Structure

- [x] **Unit Tests**: Business logic and utilities (`_test.go`)
- [ ] **Integration Tests**: Database operations and API handlers
- [x] **Benchmark Tests**: Performance critical functions

### Coverage Areas

- ✅ JWT token generation and validation
- ✅ Password hashing and verification
- ✅ UUID validation
- ✅ String utilities (snake_case conversion)
- ✅ Pointer helpers
- ✅ Repository operations (CRUD)
- ✅ Service layer logic
- ✅ HTTP handlers
- ✅ Event publishing/consuming

## Security Features

- **JWT Authentication**: Self-signed tokens with expiration
- **Password Hashing**: bcrypt with salt
- **Input Validation**: go-playground/validator
- **SQL Injection Protection**: Prepared statements
- **CORS**: Configurable cross-origin policies
- **Audit Logging**: Complete audit trail
- **Activity Tracking**: User action monitoring

## Event-Driven Architecture

The system uses NATS for asynchronous event processing:

### Published Events

- `order.created` - New order submitted
- `order.status_changed` - Order status update
- `drone.location_updated` - Drone location change
- `drone.status_changed` - Drone status update

### Event Consumers

- Activity log writer
- Audit log writer
- SMS notification sender
- Email notification sender (future)

## Monitoring & Observability

### Health Checks

```bash
curl http://localhost:8080/health
```

### NATS Monitoring

- URL: <http://localhost:8222>
- Endpoints: `/varz`, `/connz`, `/routez`, `/subsz`

### Redis Monitoring

```bash
docker exec -it redis redis-cli MONITOR
```

### Database Queries

Access PgAdmin at <http://localhost:5050> for query analysis and performance tuning.

## Environment Variables

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=postgres

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# NATS
NATS_URL=nats://localhost:4222
NATS_SERVERS=nats://localhost:4222

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
```

## Project Status

### Completed Features 

- JWT authentication system
- User management (admin, enduser, drone)
- Order CRUD operations
- Drone fleet management
- Real-time location tracking
- Status workflow management
- Order reservation and pickup
- Order withdrawal
- Drone handoff on failure
- Bulk order retrieval
- Database migrations
- Docker containerization
- Event-driven architecture

### Pending Features 

- Order delivery/failure marking
- ETA calculation algorithm
- SMS notifications (via NATS)
- Email notifications
- WebSocket support for real-time updates
- API rate limiting
- Prometheus metrics
- API documentation (Swagger)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices and idioms
- Run `go fmt` before committing
- Ensure all tests pass
- Add tests for new features
- Update documentation

## Authors

- Amer mohammmed - Initial work

## Acknowledgments

- Clean Architecture by Robert C. Martin
- Go community for excellent libraries
- NATS.io for messaging infrastructure
