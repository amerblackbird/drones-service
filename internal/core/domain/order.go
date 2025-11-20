package domain

type OrderStatus string

// Based on the order status constants, here's the difference between `picked_up` and `in_transit`:
// **`picked_up`**: The drone has successfully picked up the package from the origin location but hasn't started moving toward the destination yet.
// **`in_transit`**: The drone is actively traveling with the package from the origin to the destination location.
// The flow would typically be: `pending` → `reserved` → `picked_up` → `in_transit` → `arrived` → `delivered`
// The flow would typically be: `handoff` → `reassigned` → `in_transit`→ `arrived` → `delivered`
// So `picked_up` is the moment of collection, while `in_transit` indicates active delivery movement.

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusReserved   OrderStatus = "reserved"
	OrderStatusPickedUp   OrderStatus = "picked_up"
	OrderStatusInTransit  OrderStatus = "in_transit"
	OrderStatusArrived    OrderStatus = "arrived"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusHandoff    OrderStatus = "handoff"
	OrderStatusReassigned OrderStatus = "reassigned"
)

type Order struct {
	BaseModel
	OrderNumber          string      `json:"order_number" gorm:"uniqueIndex"`
	UserID               string      `json:"user_id"`
	ReceiverName         *string     `json:"receiver_name,omitempty"`
	ReceiverPhone        *string     `json:"receiver_phone,omitempty"`
	DeliveryNote         *string     `json:"delivery_note,omitempty"`
	PackageWeightKg      *float64    `json:"package_weight_kg,omitempty"`
	OriginAddress        string      `json:"origin_address"`
	OriginLat            float64     `json:"origin_lat"`
	OriginLon            float64     `json:"origin_lon"`
	DestinationAddress   string      `json:"destination_address"`
	DestinationLat       float64     `json:"destination_lat"`
	DestinationLon       float64     `json:"destination_lon"`
	Status               OrderStatus `json:"status" gorm:"default:pending"`
	ScheduledAt          *string     `json:"scheduled_at,omitempty"`
	DeliveredAt          *string     `json:"delivered_at,omitempty"`
	CancelledAt          *string     `json:"cancelled_at,omitempty"`
	DroneID              *string     `json:"drone_id"`
	DeliveredByDroneID   *string     `json:"delivered_by_drone_id,omitempty"`
	WithdrawnAt          *string     `json:"withdrawn_at,omitempty"`
	CurrentLat           *float64    `json:"current_lat,omitempty"`
	CurrentLon           *float64    `json:"current_lon,omitempty"`
	CurrentAltitude      *float64    `json:"current_altitude,omitempty"`
	LastLocationUpdateAt *string     `json:"last_location_update_at,omitempty"`
	EstimatedArrivalAt   *string     `json:"estimated_arrival_at,omitempty"`
}

type OrderDTO struct {
	ID                   string      `json:"id"`
	Status               OrderStatus `json:"status"`
	OrderNumber          string      `json:"order_number"`
	UserID               string      `json:"user_id"`
	ReceiverName         *string     `json:"receiver_name,omitempty"`
	ReceiverPhone        *string     `json:"receiver_phone,omitempty"`
	DeliveryNote         *string     `json:"delivery_note,omitempty"`
	PackageWeightKg      *float64    `json:"package_weight_kg,omitempty"`
	OriginAddress        string      `json:"origin_address"`
	OriginLat            float64     `json:"origin_lat"`
	OriginLon            float64     `json:"origin_lon"`
	DestinationAddress   string      `json:"destination_address"`
	DestinationLat       float64     `json:"destination_lat"`
	DestinationLon       float64     `json:"destination_lon"`
	ScheduledAt          *string     `json:"scheduled_at"`
	DeliveredAt          *string     `json:"delivered_at"`
	CancelledAt          *string     `json:"cancelled_at"`
	DroneID              *string     `json:"drone_id"`
	DeliveredByDroneID   *string     `json:"delivered_by_drone_id"`
	CreatedAt            string      `json:"created_at"`
	UpdatedAt            string      `json:"updated_at"`
	Active               bool        `json:"active"`
	CreatedByID          *string     `json:"created_by_id"`
	UpdatedByID          *string     `json:"updated_by_id"`
	WithdrawnAt          *string     `json:"withdrawn_at"`
	CurrentLat           *float64    `json:"current_lat"`
	CurrentLon           *float64    `json:"current_lon"`
	CurrentAltitude      *float64    `json:"current_altitude"`
	LastLocationUpdateAt *string     `json:"last_location_update_at"`
	EstimatedArrivalAt   *string     `json:"estimated_arrival_at"`
}
type CreateOrderRequest struct {
	ReceiverName       *string  `json:"receiver_name" validate:"omitempty,min=1"`
	ReceiverPhone      *string  `json:"receiver_phone" validate:"saudiphonenumber,min=10"`
	DeliveryNote       *string  `json:"delivery_note" validate:"omitempty,max=255"`
	PackageWeightKg    *float64 `json:"package_weight_kg,omitempty" validate:"omitempty,gt=0,lte=100"`
	OriginAddress      string   `json:"origin_address" validate:"required,min=1"`
	OriginLat          float64  `json:"origin_lat" validate:"required,saudilat"`
	OriginLon          float64  `json:"origin_lon" validate:"required,saudilon"`
	DestinationAddress string   `json:"destination_address" validate:"required,min=1"`
	DestinationLat     float64  `json:"destination_lat" validate:"required,saudilat,nefield=OriginLat"`
	DestinationLon     float64  `json:"destination_lon" validate:"required,saudilon"`
	ScheduledAt        *string  `json:"scheduled_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateOrderRequest struct {
	ReceiverName         *string      `json:"receiver_name,omitempty"`
	ReceiverPhone        *string      `json:"receiver_phone,omitempty"`
	DeliveryNote         *string      `json:"delivery_note,omitempty"`
	PackageWeightKg      *float64     `json:"package_weight_kg,omitempty"`
	OriginAddress        *string      `json:"origin_address,omitempty"`
	OriginLat            *float64     `json:"origin_lat,omitempty"`
	OriginLon            *float64     `json:"origin_lon,omitempty"`
	DestinationAddress   *string      `json:"destination_address,omitempty"`
	DestinationLat       *float64     `json:"destination_lat,omitempty"`
	DestinationLon       *float64     `json:"destination_lon,omitempty"`
	Status               *OrderStatus `json:"status,omitempty"`
	ScheduledAt          *string      `json:"scheduled_at,omitempty"`
	DeliveredAt          *string      `json:"delivered_at,omitempty"`
	CancelledAt          *string      `json:"cancelled_at,omitempty"`
	DeliveredByDroneID   *string      `json:"delivered_by_drone_id,omitempty"`
	FailedAt             *string      `json:"failed_at,omitempty"`
	DroneID              *string      `json:"drone_id,omitempty"`
	WithdrawnAt          *string      `json:"withdrawn_at,omitempty"`
	CurrentLat           *float64     `json:"current_lat,omitempty"`
	CurrentLon           *float64     `json:"current_lon,omitempty"`
	CurrentAltitude      *float64     `json:"current_altitude,omitempty"`
	LastLocationUpdateAt *string      `json:"last_location_update_at,omitempty"`
	EstimatedArrivalAt   *string      `json:"estimated_arrival_at,omitempty"`
	UpdatedByID          *string      `json:"updated_by_id"`
}

type UpdateStatusRequest struct {
	DroneID     string
	UpdatedByID string
	Status      OrderStatus
	CancelledAt *string
	DeliveredAt *string `json:"delivered_at,omitempty"`
	FailAt      *string `json:"fail_at,omitempty"`
	WithdrawnAt *string `json:"withdrawn_at,omitempty"`
}

type UpdateOrderLocationRequest struct {
	Lat      float64 `json:"lat" validate:"required,saudilat"`
	Lng      float64 `json:"lng" validate:"required,saudilon"`
	Alti     float64 `json:"alti" validate:"required,gte=0,lte=5000"`
	SpeedKmh float64 `json:"speed_kmh" validate:"required,gte=0"`
}

type OrderFilter struct {
	Status             *OrderStatus `json:"status,omitempty"`
	UserID             *string      `json:"user_id,omitempty"`
	Active             *bool        `json:"active,omitempty"`
	DroneID            *string      `json:"drone_id,omitempty"`
	DeliveredByDroneID *string      `json:"delivered_by_drone_id,omitempty"`
	DestinationAddress *string      `json:"destination_address,omitempty"`
	CreatedAtFrom      *string      `json:"created_at_from,omitempty"`
	CreatedAtTo        *string      `json:"created_at_to,omitempty"`
	ScheduledAtFrom    *string      `json:"scheduled_at_from,omitempty"`
	ScheduledAtTo      *string      `json:"scheduled_at_to,omitempty"`
	MinWeight          *float64     `json:"min_weight,omitempty"`
	MaxWeight          *float64     `json:"max_weight,omitempty"`
	ReceiverPhone      *string      `json:"receiver_phone,omitempty"`
	ReceiverName       *string      `json:"receiver_name,omitempty"`
	OriginAddress      *string      `json:"origin_address,omitempty"`
}

func (o *Order) ToDTO() *OrderDTO {
	return &OrderDTO{
		ID:                   o.ID,
		Status:               o.Status,
		OrderNumber:          o.OrderNumber,
		UserID:               o.UserID,
		ReceiverName:         o.ReceiverName,
		ReceiverPhone:        o.ReceiverPhone,
		DeliveryNote:         o.DeliveryNote,
		PackageWeightKg:      o.PackageWeightKg,
		OriginAddress:        o.OriginAddress,
		OriginLat:            o.OriginLat,
		OriginLon:            o.OriginLon,
		DestinationAddress:   o.DestinationAddress,
		DestinationLat:       o.DestinationLat,
		DestinationLon:       o.DestinationLon,
		ScheduledAt:          o.ScheduledAt,
		DeliveredAt:          o.DeliveredAt,
		CancelledAt:          o.CancelledAt,
		DeliveredByDroneID:   o.DeliveredByDroneID,
		DroneID:              o.DroneID,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
		Active:               o.Active,
		WithdrawnAt:          o.WithdrawnAt,
		CurrentLat:           o.CurrentLat,
		CurrentLon:           o.CurrentLon,
		CurrentAltitude:      o.CurrentAltitude,
		LastLocationUpdateAt: o.LastLocationUpdateAt,
		EstimatedArrivalAt:   o.EstimatedArrivalAt,
	}
}

func (filter OrderFilter) IsEmpty() bool {
	return filter.Status == nil &&
		filter.UserID == nil &&
		filter.Active == nil &&
		filter.DroneID == nil &&
		filter.DeliveredByDroneID == nil &&
		filter.DestinationAddress == nil &&
		filter.CreatedAtFrom == nil &&
		filter.CreatedAtTo == nil &&
		filter.ScheduledAtFrom == nil &&
		filter.ScheduledAtTo == nil &&
		filter.MinWeight == nil &&
		filter.MaxWeight == nil &&
		filter.ReceiverPhone == nil &&
		filter.ReceiverName == nil &&
		filter.OriginAddress == nil
}

func (order *Order) IsReserved() bool {
	return order.Status == OrderStatusReserved || order.Status == OrderStatusPickedUp || order.Status == OrderStatusInTransit || order.Status == OrderStatusArrived || order.Status == OrderStatusDelivered
}
