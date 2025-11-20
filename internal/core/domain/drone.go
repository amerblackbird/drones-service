package domain

type DroneStatus string

// DroneStatus represents the current operational state of a drone in the delivery system.
//
// Workflow transitions:
//   - idle: The drone is available and ready for new delivery assignments
//   - loading: Drone transitions from idle when assigned a package and is being loaded
//   - delivering: Drone moves from loading to delivering once package is secured and delivery starts
//   - returning: After successful delivery, drone enters returning state to go back to base
//   - charging: From idle or returning, drone can enter charging when battery is low
//   - broken: Drone can transition to broken from any operational state due to malfunction
//   - under_repair: Broken drones move to under_repair when maintenance begins
//   - maintenanced: Scheduled maintenance state, typically from idle for preventive care
//
// Valid transitions:
//
//	idle -> loading, charging, maintenanced, broken
//	loading -> delivering, broken
//	delivering -> returning, broken
//	returning -> idle, charging, broken
//	charging -> idle, broken, returning
//	broken -> under_repair
//	under_repair -> maintenanced
//	maintenanced -> idle, returning
const (
	DroneStatusIdle         DroneStatus = "idle"         // Available for assignments
	DroneStatusLoading      DroneStatus = "loading"      // Being loaded with a package
	DroneStatusDelivering   DroneStatus = "delivering"   // Currently delivering a package
	DroneStatusReturning    DroneStatus = "returing"     // Returning to base after delivery
	DroneStatusCharging     DroneStatus = "charging"     // Charging its battery
	DroneStatusBroken       DroneStatus = "broken"       // Drone is broken
	DroneStatusUnderRepair  DroneStatus = "under_repair" // Undergoing repairs
	DroneStatusMaintenanced DroneStatus = "maintenanced" // Undergoing maintenance
)

type Drone struct {
	BaseModel
	DroneIdentifier      string      `json:"drone_identifier"`
	UserID               string      `json:"user_id"`
	SerialNumber         string      `json:"serial_number"`
	Model                string      `json:"model"`
	Manufacturer         string      `json:"manufacturer"`
	MaxWeightKg          float64     `json:"max_weight_kg"`
	MaxSpeedKmh          float64     `json:"max_speed_kmh"`
	MaxRangeKm           float64     `json:"max_range_km"`
	BatteryCapacityMah   int         `json:"battery_capacity_mah"`
	Status               DroneStatus `json:"status"`
	BatteryLevelPercent  *float64    `json:"battery_level_percent,omitempty"`
	CurrentLat           *float64    `json:"current_lat,omitempty"`
	CurrentLon           *float64    `json:"current_lon,omitempty"`
	CurrentAltitude      *float64    `json:"current_altitude,omitempty"`
	LastLocationUpdateAt *string     `json:"last_location_update_at,omitempty"`
	TotalFlightHours     float64     `json:"total_flight_hours"`
	TotalDeliveries      int         `json:"total_deliveries"`
	LastMaintenanceAt    *string     `json:"last_maintenance_at,omitempty"`
	NextMaintenanceDueAt *string     `json:"next_maintenance_due_at,omitempty"`
}

type DroneDTO struct {
	ID                  string      `json:"id"`
	UserID              string      `json:"user_id"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	Active              bool        `json:"active"`
	CreatedByID         *string     `json:"created_by_id"`
	UpdatedByID         *string     `json:"updated_by_id"`
	DroneIdentifier     string      `json:"drone_identifier"`
	Model               string      `json:"model"`
	SerialNumber        string      `json:"serial_number"`
	BatteryCapacity     int         `json:"battery_capacity"`
	PayloadCapacity     float64     `json:"payload_capacity"`
	Manufacturer        string      `json:"manufacturer"`
	LastChargedAt       *string     `json:"last_charged_at"`
	IsCharging          *bool       `json:"is_charging"`
	LastKnownLat        *float64    `json:"last_known_lat"`
	LastKnownLng        *float64    `json:"last_known_lng"`
	LastAltitudeM       *float64    `json:"last_altitude_m"`
	LastSpeedKmh        *float64    `json:"last_speed_kmh"`
	CurrentOrderID      *string     `json:"current_order_id"`
	CrashesCount        *int        `json:"crashes_count"`
	MaintenanceRequired *bool       `json:"maintenance_required"`
	LastMaintenanceAt   *string     `json:"last_maintenance_at"`
	NextMaintenanceAt   *string     `json:"next_maintenance_at"`
	Status              DroneStatus `json:"status"`
}

type CreateDroneRequest struct {
	Model           string  `json:"model" validate:"required,min=2,max=50"`
	SerialNumber    string  `json:"serial_number" validate:"required,alphanum,min=5,max=100"`
	Manufacturer    string  `json:"manufacturer" validate:"required,min=2,max=50"`
	BatteryCapacity int     `json:"battery_capacity" validate:"required,min=1000,max=100000"`
	PayloadCapacity float64 `json:"payload_capacity" validate:"required,min=0.1,max=500"`
	CreatedByID     string  `json:"created_by_id" validate:"required,uuid4"`
}

type UpdateDroneRequest struct {
	Model               *string      `json:"model,omitempty" validate:"omitempty,min=2,max=50"`
	SerialNumber        *string      `json:"serial_number,omitempty" validate:"omitempty,alphanum,min=5,max=100"`
	Manufacturer        *string      `json:"manufacturer,omitempty" validate:"omitempty,min=2,max=50"`
	BatteryCapacity     *int         `json:"battery_capacity,omitempty" validate:"omitempty,min=1000,max=100000"`
	PayloadCapacity     *float64     `json:"payload_capacity,omitempty" validate:"omitempty,min=0.1,max=500"`
	UpdatedByID         *string      `json:"updated_by_id,omitempty" validate:"omitempty,required,uuid4"`
	LastChargedAt       *string      `json:"last_charged_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	IsCharging          *bool        `json:"is_charging,omitempty"`
	LastKnownLat        *float64     `json:"last_known_lat,omitempty" validate:"omitempty,min=-90,max=90"`
	LastKnownLng        *float64     `json:"last_known_lng,omitempty" validate:"omitempty,min=-180,max=180"`
	LastAltitudeM       *float64     `json:"last_altitude_m,omitempty" validate:"omitempty,min=0,max=10000"`
	LastSpeedKmh        *float64     `json:"last_speed_kmh,omitempty" validate:"omitempty,min=0,max=500"`
	CurrentOrderID      *string      `json:"current_order_id,omitempty" validate:"omitempty,uuid4"`
	CrashesCount        *int         `json:"crashes_count,omitempty" validate:"omitempty,min=0"`
	MaintenanceRequired *bool        `json:"maintenance_required,omitempty"`
	LastMaintenanceAt   *string      `json:"last_maintenance_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	NextMaintenanceAt   *string      `json:"next_maintenance_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Status              *DroneStatus `json:"status,omitempty" validate:"omitempty,oneof=idle loading delivering returning charging maintenance"`
}

type DroneFilter struct {
	Status          *string  `json:"status,omitempty"`
	Statuses        []string `json:"statuses,omitempty"`
	Active          *bool    `json:"active,omitempty"`
	UserID          *string  `json:"user_id,omitempty"`
	SerialNumber    *string  `json:"serial_number,omitempty"`
	DroneIdentifier *string  `json:"drone_identifier,omitempty"`
	Lat             *float64 `json:"lat,omitempty"`
	Lon             *float64 `json:"lon,omitempty"`
	Radius          *float64 `json:"radius,omitempty"`
}

func (d *Drone) ToDTO() *DroneDTO {
	return &DroneDTO{
		ID:                d.ID,
		UserID:            d.UserID,
		Status:            d.Status,
		DroneIdentifier:   d.DroneIdentifier,
		Model:             d.Model,
		SerialNumber:      d.SerialNumber,
		BatteryCapacity:   d.BatteryCapacityMah,
		PayloadCapacity:   d.MaxWeightKg,
		Manufacturer:      d.Manufacturer,
		LastChargedAt:     d.LastLocationUpdateAt,
		LastKnownLat:      d.CurrentLat,
		LastKnownLng:      d.CurrentLon,
		LastAltitudeM:     d.CurrentAltitude,
		LastMaintenanceAt: d.LastMaintenanceAt,
		NextMaintenanceAt: d.NextMaintenanceDueAt,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
		Active:            d.Active,
		CreatedByID:       d.CreatedByID,
		UpdatedByID:       d.UpdatedByID,
	}
}

func (status DroneStatus) GetErr() error {
	switch status {
	case DroneStatusLoading:
		return ErrDroneIsLoading
	case DroneStatusDelivering:
		return ErrDroneIsDelivering
	case DroneStatusReturning:
		return ErrDroneIsReturning
	case DroneStatusCharging:
		return ErrDroneIsCharging
	case DroneStatusMaintenanced:
		return ErrDroneInMaintenance
	case DroneStatusBroken:
		return ErrDroneIsBroken
	case DroneStatusUnderRepair:
		return ErrDroneUnderRepair
	default:
		return ErrDroneMustBeIdle
	}
}

// Force workflow to update drone status to broken
//
//	idle -> loading, charging, maintenanced, broken
//	loading -> delivering, broken
//	delivering -> returning, broken
//	returning -> idle, charging, broken
//	charging -> idle, broken, returning
//	broken -> under_repair
//	under_repair -> maintenanced
//	maintenanced -> idle, returning
func (status DroneStatus) IsTransitionAllowed(from DroneStatus) bool {
	switch from {
	case DroneStatusIdle:
		return status == DroneStatusLoading || status == DroneStatusCharging || status == DroneStatusMaintenanced || status == DroneStatusBroken
	case DroneStatusLoading:
		return status == DroneStatusDelivering || status == DroneStatusBroken
	case DroneStatusDelivering:
		return status == DroneStatusReturning || status == DroneStatusBroken
	case DroneStatusReturning:
		return status == DroneStatusIdle || status == DroneStatusCharging || status == DroneStatusBroken
	case DroneStatusCharging:
		return status == DroneStatusIdle || status == DroneStatusReturning || status == DroneStatusBroken
	case DroneStatusBroken:
		return status == DroneStatusUnderRepair
	case DroneStatusUnderRepair:
		return status == DroneStatusMaintenanced
	case DroneStatusMaintenanced:
		return status == DroneStatusIdle || status == DroneStatusReturning
	default:
		return false
	}
}

func (status DroneStatus) TransitionErr() error {
	switch status {
	case DroneStatusIdle:
		return ErrIdleTransition
	case DroneStatusLoading:
		return ErrLoadingTransition
	case DroneStatusDelivering:
		return ErrDeliveringTransition
	case DroneStatusReturning:
		return ErrReturningTransition
	case DroneStatusCharging:
		return ErrChargingTransition
	case DroneStatusBroken:
		return ErrBrokenTransition
	case DroneStatusUnderRepair:
		return ErrUnderRepairTransition
	case DroneStatusMaintenanced:
		return ErrMaintenancedTransition
	default:
		return nil
	}
}
