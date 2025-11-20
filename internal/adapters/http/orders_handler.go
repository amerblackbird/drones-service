package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"drones/internal/core/domain"
	"drones/internal/ports"
	"drones/pkg/utils"
)

type OrdersHandler struct {
	service        ports.OrdersService
	validator      *validator.Validate
	logger         ports.Logger
	eventPublisher ports.EventPublisher
}

func NewOrdersHandler(service ports.OrdersService, eventPublisher ports.EventPublisher, logger ports.Logger) *OrdersHandler {
	return &OrdersHandler{
		service:        service,
		validator:      domain.NewValidator(),
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}

// RegisterRoutes registers all order routes
func (h *OrdersHandler) RegisterRoutes(r *mux.Router) {
	// Get current order
	r.Handle("/current", DroneGuard(http.HandlerFunc(h.HandleGetCurrentOrder))).Methods("GET")

	r.HandleFunc("", h.HandleCreateOrder).Methods("POST")
	r.HandleFunc("", h.HandleListOrders).Methods("GET")
	r.HandleFunc("/{id}", h.HandleGetOrder).Methods("GET")
	r.Handle("/{id}", AdminGuard(http.HandlerFunc(h.HandleUpdateOrder))).Methods("PUT")

	r.HandleFunc("/{id}", h.HandleUpdateOrder).Methods("PUT")

	// r.HandleFunc("/{id}", h.HandleDeleteOrder).Methods("DELETE")
	r.Handle("/{id}/withdraw", EndUserGuard(http.HandlerFunc(h.HandleOrderWithdrawn))).Methods("POST")

	// Drone actions
	r.Handle("/{id}/reserve", DroneGuard(http.HandlerFunc(h.HandleReserveOrder))).Methods("POST")
	r.Handle("/{id}/confirm-pickup", DroneGuard(http.HandlerFunc(h.HandleConfirmPickup))).Methods("POST")
	r.Handle("/{id}/start-transit", DroneGuard(http.HandlerFunc(h.HandleStartTransit))).Methods("POST")
	r.Handle("/{id}/confirm-arrival", DroneGuard(http.HandlerFunc(h.HandleConfirmArrival))).Methods("POST")
	r.Handle("/{id}/confirm-delivery", DroneGuard(http.HandlerFunc(h.HandleConfirmDelivery))).Methods("POST")
	r.Handle("/{id}/delivery-failed", DroneGuard(http.HandlerFunc(h.HandleDeliveryFailed))).Methods("POST")

	// Handoff endpoint
	// r.HandleFunc("/{id}/handoff", h.HandleOrderHandoff).Methods("POST")
	// r.HandleFunc("/{id}/reassign", h.HandleReassign).Methods("POST")

	// Location
	// r.HandleFunc("/{id}/heartbeat", h.HandleHeartbeat).Methods("POST")

}

// HandleCreateOrder creates a new order
func (h *OrdersHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	meta := domain.ExtractRequestInfo(r)
	h.logger.Info("Creating order", zap.Any("data", meta))

	var request domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid request body", err))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		ResponseWithValidationError(w, http.StatusBadRequest, domain.GetValidationErrors(err.(validator.ValidationErrors)))
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

	order, err := h.service.CreateOrder(r.Context(), user.ID, &request)
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusCreated, order.ToDTO())
}

// HandleGetOrder retrieves a order by ID
func (h *OrdersHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Orders ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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

	// Parse filter parameters
	filter := domain.OrderFilter{}
	switch user.Type {
	case domain.UserTypeEnduser:
		filter.UserID = &user.ID
	case domain.UserTypeDrone:
		filter.DroneID = user.DroneId
		filter.DeliveredByDroneID = user.DroneId
	case domain.UserTypeAdmin:
		// Admin can see all orders
	}

	order, err := h.service.GetOrderByID(r.Context(), id, filter)
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

// HandleUpdateOrder updates an existing order
func (h *OrdersHandler) HandleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Orders ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
		return
	}
	user, ok := UserFromContext(r.Context())
	if !ok || user == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UnauthorizedError,
			Message: "User not found in context",
		})
		return
	}

	// Reject if user is of type drone
	if user.Type == domain.UserTypeDrone {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UnauthorizedError,
			Message: "Only drone users can update orders",
		})
	}

	var request domain.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithValidationError(w, http.StatusBadRequest, domain.GetValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		errors := domain.GetValidationErrors(err.(validator.ValidationErrors))
		ResponseWithValidationError(w, http.StatusBadRequest, errors)
		return
	}

	order, err := h.service.UpdateOrder(r.Context(), id, &request, domain.OrderFilter{})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

// HandleDeleteOrder soft deletes a order
func (h *OrdersHandler) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ResponseWithResouseNotFound(w, "Orders ID")
		return
	}

	if !utils.ValidateUUID(id) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
		return
	}

	deletedBy := r.Header.Get("X-User-ID")
	if deletedBy == "" {
		deletedBy = "system"
	}

	err := h.service.DeleteOrder(r.Context(), id, domain.OrderFilter{})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusNoContent, map[string]interface{}{})
}

// HandleListOrders retrieves orders with filtering and pagination
func (h *OrdersHandler) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := GetPaginationParams(r)
	if err != nil {
		ResponseWithError(w, err)
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

	// Parse filter parameters
	filter := domain.OrderFilter{}
	switch user.Type {
	case domain.UserTypeEnduser:
		filter.UserID = &user.ID
	case domain.UserTypeDrone:
		filter.DroneID = user.DroneId
		filter.DeliveredByDroneID = user.DroneId
	case domain.UserTypeAdmin:
		// Admin can see all orders
	}

	if status := r.URL.Query().Get("status"); status != "" {
		orderStatus := domain.OrderStatus(status)
		filter.Status = &orderStatus
	}

	if active := r.URL.Query().Get("active"); active != "" {
		if activeBool, err := strconv.ParseBool(active); err == nil {
			filter.Active = &activeBool
		}
	}

	options := domain.PaginationOption[domain.OrderFilter]{
		Filter: &filter,
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.service.ListOrders(r.Context(), options)
	if err != nil {
		ResponseWithError(w, err)
		return
	}
	ResponseWithJSON(w, http.StatusOK, result)
}

func (s *OrdersHandler) HandleOrderWithdrawn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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

	order, err := s.service.Withdraw(r.Context(), orderID, user.ID, domain.OrderFilter{
		UserID: &user.ID,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleReserveOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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

	order, err := h.service.Reserve(r.Context(), orderID, user.ID, domain.OrderFilter{})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleConfirmPickup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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
	h.logger.Info("User Drone ID", zap.Any("drone_id", user.DroneId))
	if user.DroneId == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "Drone ID not found for user",
		})
		return
	}

	order, err := h.service.ConfirmPickup(r.Context(), orderID, user.ID, domain.OrderFilter{
		DroneID: user.DroneId,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleStartTransit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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
	if user.DroneId == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "Drone ID not found for user",
		})
		return
	}

	order, err := h.service.StartTransit(r.Context(), orderID, user.ID, domain.OrderFilter{
		DroneID: user.DroneId,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleConfirmArrival(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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
	if user.DroneId == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "Drone ID not found for user",
		})
		return
	}
	order, err := h.service.ConfirmArrived(r.Context(), orderID, user.ID, domain.OrderFilter{
		DroneID: user.DroneId,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleConfirmDelivery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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
	if user.DroneId == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "Drone ID not found for user",
		})
		return
	}

	order, err := h.service.ConfirmDelivery(r.Context(), orderID, user.ID, domain.OrderFilter{
		DroneID: user.DroneId,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleDeliveryFailed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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
	if user.DroneId == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "Drone ID not found for user",
		})
		return
	}

	order, err := h.service.DeliveryFailed(r.Context(), orderID, user.ID, domain.OrderFilter{
		DroneID: user.DroneId,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleReassign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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

	order, err := h.service.Reassign(r.Context(), orderID, user.ID, domain.OrderFilter{})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		ResponseWithResouseNotFound(w, "Order ID")
		return
	}

	if !utils.ValidateUUID(orderID) {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid order ID format", nil))
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

	var request domain.UpdateOrderLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid request body", err))
		return
	}

	order, err := h.service.UpadateOrderLocation(r.Context(), user.ID, orderID, request.Lat, request.Lng, request.Alti, domain.OrderFilter{
		DroneID: &user.ID,
	})
	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}

func (h *OrdersHandler) HandleGetCurrentOrder(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok || user == nil {
		ResponseWithCustomError(w, http.StatusUnauthorized, domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "User not found in context",
		})
		return
	}

	order, err := h.service.GetOrderByFilter(r.Context(), domain.OrderFilter{
		DroneID: &user.ID,
	})

	if err != nil {
		ResponseWithError(w, err)
		return
	}

	ResponseWithJSON(w, http.StatusOK, order.ToDTO())
}
