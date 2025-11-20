package domain

type HeartbeatRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,saudilat"`
	Longitude float64 `json:"longitude" validate:"required,saudilon"`
	Altitude  float64 `json:"altitude" validate:"required,gte=0"`
	Battery   int     `json:"battery" validate:"required,gte=0,lte=100"`
}
