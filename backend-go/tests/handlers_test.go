package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/handlers"
	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
)

func TestHandlersHTTP(t *testing.T) {
	jwtSecret := "capstone-test-secret-key-987"
	store := db.NewMemoryStore()
	db.SeedMemoryStore(store)

	authH := handlers.NewAuthHandler(store, jwtSecret)
	medH := handlers.NewMedicinesHandler(store)
	adminH := handlers.NewAdminHandler(store)

	t.Run("Auth Register creates new user and returns JWT token", func(t *testing.T) {
		body := `{"name":"Bob Builder","email":"bob@pharmacy.local","password":"SecretPassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		authH.Register(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d. Body: %s", rec.Code, rec.Body.String())
		}

		var res models.AuthResponse
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("failed decoding json: %v", err)
		}

		if res.Token == "" || res.User.Email != "bob@pharmacy.local" || res.User.Role != models.RoleCustomer {
			t.Errorf("unexpected register response: %+v", res)
		}
	})

	t.Run("Auth Login returns token for seeded demo customer", func(t *testing.T) {
		body := `{"email":"customer@pharmacy.local","password":"CustomerPass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		authH.Login(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d. Body: %s", rec.Code, rec.Body.String())
		}

		var res models.AuthResponse
		_ = json.NewDecoder(rec.Body).Decode(&res)
		if res.Token == "" || res.User.Role != models.RoleCustomer {
			t.Errorf("unexpected login response: %+v", res)
		}
	})

	t.Run("Medicines catalog search filters by term case-insensitively", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/medicines?q=amox&page=1&limit=10", nil)
		rec := httptest.NewRecorder()

		medH.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}

		var listRes handlers.MedicinesListResponse
		if err := json.NewDecoder(rec.Body).Decode(&listRes); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(listRes.Data) == 0 {
			t.Fatal("expected search for 'amox' to find items, got 0")
		}
		for _, item := range listRes.Data {
			if item.Name == "" || item.GenericName == "" {
				t.Error("medicine items must have name and generic name")
			}
		}
	})

	t.Run("Admin update role changes user role", func(t *testing.T) {
		cust, err := store.GetUserByEmail(context.Background(), "customer@pharmacy.local")
		if err != nil {
			t.Fatalf("failed to find customer: %v", err)
		}

		body := `{"role":"staff"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/users/"+cust.UserID+"/role", bytes.NewBufferString(body))
		req.SetPathValue("id", cust.UserID)
		req.Header.Set("Content-Type", "application/json")

		// Admin user context
		adminUser := &models.User{UserID: "admin-1", Role: models.RoleAdmin}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, adminUser)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		adminH.UpdateRole(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d. Body: %s", rec.Code, rec.Body.String())
		}

		updated, _ := store.GetUserByID(context.Background(), cust.UserID)
		if updated.Role != models.RoleStaff {
			t.Errorf("expected role staff, got %s", updated.Role)
		}
	})
}
