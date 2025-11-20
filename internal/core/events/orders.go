package events

import "drones/internal/core/domain"

type OrderCreatedEvent struct {
	OrderID            string  `json:"order_id"`
	UserID             string  `json:"user_id"`
	OriginAddress      string  `json:"origin_address"`
	OriginLat          float64 `json:"origin_lat"`
	OriginLon          float64 `json:"origin_lon"`
	DestinationAddress string  `json:"destination_address"`
	DestinationLat     float64 `json:"destination_lat"`
	DestinationLon     float64 `json:"destination_lon"`
}

type OrderUpdatedEvent struct {
	OrderID         string             `json:"order_id"`
	UserID          string             `json:"user_id"`
	DroneID         string             `json:"drone_id,omitempty"`
	Status          domain.OrderStatus `json:"status"`
	CurrentLat      *float64           `json:"current_lat,omitempty"`
	CurrentLon      *float64           `json:"current_lon,omitempty"`
	CurrentAltitude *float64           `json:"current_altitude,omitempty"`
}

type OrderWithdrawnEvent struct {
	OrderID string             `json:"order_id"`
	UserID  string             `json:"user_id"`
	Status  domain.OrderStatus `json:"status"`
}

type OrderReservedEvent struct {
	OrderID string             `json:"order_id"`
	DroneID string             `json:"drone_id"`
	Status  domain.OrderStatus `json:"status"`
}
