package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/models"
)

type AdminHandler struct {
	store db.Store
}

func NewAdminHandler(store db.Store) *AdminHandler {
	return &AdminHandler{store: store}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}
	JSONResponse(w, http.StatusOK, users)
}

func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	var req models.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	targetRole := models.UserRole(strings.ToLower(strings.TrimSpace(string(req.Role))))
	switch targetRole {
	case models.RoleCustomer, models.RoleStaff, models.RoleAdmin:
	default:
		JSONError(w, http.StatusBadRequest, "Role must be 'customer', 'staff', or 'admin'")
		return
	}

	if err := h.store.UpdateUserRole(r.Context(), id, targetRole); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "User not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	user, err := h.store.GetUserByID(r.Context(), id)
	if err != nil {
		JSONResponse(w, http.StatusOK, map[string]string{"message": "Role updated successfully"})
		return
	}

	JSONResponse(w, http.StatusOK, user)
}
