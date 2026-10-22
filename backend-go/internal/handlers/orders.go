package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
)

type OrdersHandler struct {
	store db.Store
}

func NewOrdersHandler(store db.Store) *OrdersHandler {
	return &OrdersHandler{store: store}
}

func (h *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Items) == 0 {
		JSONError(w, http.StatusBadRequest, "Order must contain at least one medicine item")
		return
	}

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			JSONError(w, http.StatusBadRequest, "Quantity must be greater than zero for all items")
			return
		}
	}

	order, err := h.store.CreateOrderWithItems(r.Context(), authUser.UserID, req.PrescriptionID, req.Items)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrInsufficientStock):
			JSONError(w, http.StatusConflict, err.Error())
		case errors.Is(err, db.ErrRxPrescriptionReq):
			JSONError(w, http.StatusBadRequest, "Prescription required: one or more items require an approved prescription owned by your account")
		case errors.Is(err, db.ErrPrescriptionNotOwn):
			JSONError(w, http.StatusBadRequest, "Invalid prescription: prescription not found or not owned by you")
		default:
			JSONError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	JSONResponse(w, http.StatusCreated, order)
}

func (h *OrdersHandler) Mine(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orders, err := h.store.ListOrdersByUser(r.Context(), authUser.UserID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve your orders")
		return
	}

	JSONResponse(w, http.StatusOK, orders)
}

func (h *OrdersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := h.store.GetOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Order not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if order.UserID != authUser.UserID && authUser.Role != models.RoleStaff && authUser.Role != models.RoleAdmin {
		JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	JSONResponse(w, http.StatusOK, order)
}

func (h *OrdersHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.store.ListAllOrders(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve orders list")
		return
	}

	JSONResponse(w, http.StatusOK, orders)
}

func (h *OrdersHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	var req models.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	targetStatus := models.OrderStatus(strings.ToUpper(strings.TrimSpace(string(req.Status))))
	updated, err := h.store.UpdateOrderStatus(r.Context(), id, targetStatus, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			JSONError(w, http.StatusNotFound, "Order not found")
		case errors.Is(err, db.ErrInvalidTransition):
			JSONError(w, http.StatusBadRequest, err.Error())
		default:
			JSONError(w, http.StatusInternalServerError, "Failed to update order status")
		}
		return
	}

	JSONResponse(w, http.StatusOK, updated)
}
