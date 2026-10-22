package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/models"
)

type MedicinesHandler struct {
	store db.Store
}

func NewMedicinesHandler(store db.Store) *MedicinesHandler {
	return &MedicinesHandler{store: store}
}

type MedicinesListResponse struct {
	Data  []models.Medicine `json:"data"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

func (h *MedicinesHandler) List(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	meds, total, err := h.store.ListMedicines(r.Context(), q, page, limit)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to query medicines catalog")
		return
	}

	JSONResponse(w, http.StatusOK, MedicinesListResponse{
		Data:  meds,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *MedicinesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Medicine ID is required")
		return
	}

	med, err := h.store.GetMedicineByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Medicine not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve medicine")
		return
	}

	JSONResponse(w, http.StatusOK, med)
}

func (h *MedicinesHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var req models.Medicine
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.GenericName = strings.TrimSpace(req.GenericName)
	if req.Name == "" || req.GenericName == "" || req.Price < 0 || req.StockQty < 0 {
		JSONError(w, http.StatusBadRequest, "Name, generic name, valid price, and non-negative stock quantity are required")
		return
	}

	req.CreatedAt = time.Now().UTC()
	if err := h.store.CreateMedicine(r.Context(), &req); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to create medicine")
		return
	}

	JSONResponse(w, http.StatusCreated, req)
}

func (h *MedicinesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Medicine ID is required")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	var req models.Medicine
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.MedicineID = id
	req.Name = strings.TrimSpace(req.Name)
	req.GenericName = strings.TrimSpace(req.GenericName)
	if req.Name == "" || req.GenericName == "" || req.Price < 0 || req.StockQty < 0 {
		JSONError(w, http.StatusBadRequest, "Name, generic name, valid price, and non-negative stock quantity are required")
		return
	}

	if err := h.store.UpdateMedicine(r.Context(), &req); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Medicine not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to update medicine")
		return
	}

	JSONResponse(w, http.StatusOK, req)
}

func (h *MedicinesHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Medicine ID is required")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	var req models.UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.StockQty != nil && *req.StockQty < 0 {
		JSONError(w, http.StatusBadRequest, "Stock quantity cannot be negative")
		return
	}
	if req.Price != nil && *req.Price < 0 {
		JSONError(w, http.StatusBadRequest, "Price cannot be negative")
		return
	}

	updated, err := h.store.UpdateMedicineStockAndPrice(r.Context(), id, req.StockQty, req.Price)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Medicine not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to update stock/price")
		return
	}

	JSONResponse(w, http.StatusOK, updated)
}
