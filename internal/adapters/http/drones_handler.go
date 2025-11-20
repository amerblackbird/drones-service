package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"drones/internal/core/domain"
	"drones/internal/ports"
	"drones/pkg/utils"
)

type DronesHandler struct {
	service        ports.DronesService
	validator      *validator.Validate
	logger         ports.Logger
	eventPublisher ports.EventPublisher
}

func NewDronesHandler(service ports.DronesService, eventPublisher ports.EventPublisher, logger ports.Logger) *DronesHandler {
	return &DronesHandler{
		service:        service,
		validator:      domain.NewValidator(),
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}

// RegisterRoutes registers all drone routes
func (h *DronesHandler) RegisterRoutes(r *mux.Router) {
	r.Handle("", AdminGuard(http.HandlerFunc(h.HandleListDrones))).Methods("GET")
	r.Handle("/{id}", AdminGuard(http.HandlerFunc(h.HandleGetDrone))).Methods("GET")
	r.Handle("/{id}", AdminGuard(http.HandlerFunc(h.HandleUpdateDrone))).Methods("PUT")
	r.Handle("/{id}/status", AdminGuard(http.HandlerFunc(h.HandleStatusUpdated))).Methods("POST")

	// Drone submit Location
	r.Handle("/{id}/heartbeat", DroneGuard(http.HandlerFunc(h.HandleHeartbeat))).Methods("POST")
}

// HandleGetDrone retrieves a drone by ID
func (h *DronesHandler) HandleGetDrone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Drones ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid drone ID format", nil))
		return
	}

	drone, err := h.service.GetDroneByID(r.Context(), id)
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, drone.ToDTO())
}

// HandleUpdateDrone updates an existing drone
func (h *DronesHandler) HandleUpdateDrone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Drones ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid drone ID format", nil))
		return
	}

	var request domain.UpdateDroneRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithValidationError(w, http.StatusBadRequest, domain.GetValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		errors := domain.GetValidationErrors(err.(validator.ValidationErrors))
		ResponseWithValidationError(w, http.StatusBadRequest, errors)
		return
	}
	user, ok := UserFromContext(r.Context())
	if !ok || user == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "User not found in context",
		})
		return
	}

	drone, err := h.service.UpdateDrone(r.Context(), id, &domain.UpdateDroneRequest{
		Model:               request.Model,
		SerialNumber:        request.SerialNumber,
		Manufacturer:        request.Manufacturer,
		BatteryCapacity:     request.BatteryCapacity,
		PayloadCapacity:     request.PayloadCapacity,
		IsCharging:          request.IsCharging,
		LastChargedAt:       request.LastChargedAt,
		LastKnownLat:        request.LastKnownLat,
		LastKnownLng:        request.LastKnownLng,
		LastAltitudeM:       request.LastAltitudeM,
		LastSpeedKmh:        request.LastSpeedKmh,
		CurrentOrderID:      request.CurrentOrderID,
		CrashesCount:        request.CrashesCount,
		MaintenanceRequired: request.MaintenanceRequired,
		LastMaintenanceAt:   request.LastMaintenanceAt,
		NextMaintenanceAt:   request.NextMaintenanceAt,
		Status:              request.Status,
		UpdatedByID:         &user.ID,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, drone.ToDTO())
}

// HandleListDrones retrieves drones with filtering and pagination
func (h *DronesHandler) HandleListDrones(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := GetPaginationParams(r)
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	// Parse filter parameters
	filter := &domain.DroneFilter{}

	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}

	if active := r.URL.Query().Get("active"); active != "" {
		if activeBool, err := strconv.ParseBool(active); err == nil {
			filter.Active = &activeBool
		}
	}

	options := domain.PaginationOption[domain.DroneFilter]{
		Filter: filter,
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.service.ListDrones(r.Context(), options)
	if err != nil {
		ResponseWithError(w, err)
		return
	}
	ResponseWithJSON(w, http.StatusOK, result)
}

func (h *DronesHandler) HandleStatusUpdated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Drones ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid drone ID format", nil))
		return
	}

	var request struct {
		Status domain.DroneStatus `json:"status" validate:"required,oneof=idle loading delivering returning charging broken under_repair maintenanced"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithValidationError(w, http.StatusBadRequest, domain.GetValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		errors := domain.GetValidationErrors(err.(validator.ValidationErrors))
		ResponseWithValidationError(w, http.StatusBadRequest, errors)
		return
	}
	user, ok := UserFromContext(r.Context())
	if !ok || user == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "User not found in context",
		})
		return
	}

	drone, err := h.service.UpdateDroneStatus(r.Context(), user.ID, id, request.Status)
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, drone.ToDTO())
}

func (h *DronesHandler) HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Drones ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid drone ID format", nil))
		return
	}

	var request domain.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithValidationError(w, http.StatusBadRequest, domain.GetValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		errors := domain.GetValidationErrors(err.(validator.ValidationErrors))
		ResponseWithValidationError(w, http.StatusBadRequest, errors)
		return
	}
	user, ok := UserFromContext(r.Context())
	if !ok || user == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "User not found in context",
		})
		return
	}

	updatedDrone, err := h.service.ProcessHeartbeat(r.Context(), *user.DroneId, user.ID, request)
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, updatedDrone.ToDTO())
}
