package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
)

func TestAuthAndRoleMiddleware(t *testing.T) {
	secret := "test-secret-key-12345"

	customerUser := &models.User{
		UserID: "cust-01",
		Email:  "cust@test.local",
		Role:   models.RoleCustomer,
	}

	staffUser := &models.User{
		UserID: "staff-01",
		Email:  "staff@test.local",
		Role:   models.RoleStaff,
	}

	custToken, err := middleware.GenerateToken(customerUser, secret)
	if err != nil {
		t.Fatalf("failed generating token: %v", err)
	}

	staffToken, err := middleware.GenerateToken(staffUser, secret)
	if err != nil {
		t.Fatalf("failed generating staff token: %v", err)
	}

	t.Run("Valid token decodes correct claims", func(t *testing.T) {
		claims, err := middleware.ValidateToken(custToken, secret)
		if err != nil {
			t.Fatalf("expected valid token, got: %v", err)
		}
		if claims.UserID != customerUser.UserID || claims.Email != customerUser.Email || claims.Role != customerUser.Role {
			t.Errorf("claims mismatch: %+v", claims)
		}
	})

	t.Run("Invalid signature rejects token", func(t *testing.T) {
		_, err := middleware.ValidateToken(custToken, "wrong-secret-key")
		if err == nil {
			t.Error("expected error with wrong secret, got nil")
		}
	})

	t.Run("RequireAuth rejects missing token", func(t *testing.T) {
		called := false
		protectedHandler := middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		rec := httptest.NewRecorder()

		// Run through auth middleware without header
		middleware.AuthMiddleware(secret)(protectedHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d", rec.Code)
		}
		if called {
			t.Error("protected handler should not have been called")
		}
	})

	t.Run("RequireAuth allows valid token", func(t *testing.T) {
		called := false
		var userCtx *models.User
		protectedHandler := middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			called = true
			userCtx, _ = middleware.GetUserFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		req.Header.Set("Authorization", "Bearer "+custToken)
		rec := httptest.NewRecorder()

		middleware.AuthMiddleware(secret)(protectedHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", rec.Code)
		}
		if !called || userCtx == nil || userCtx.UserID != "cust-01" {
			t.Error("expected user context to be populated")
		}
	})

	t.Run("RequireRole forbids customer from staff endpoint", func(t *testing.T) {
		called := false
		staffHandler := middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/prescriptions/queue", nil)
		req.Header.Set("Authorization", "Bearer "+custToken) // customer token
		rec := httptest.NewRecorder()

		middleware.AuthMiddleware(secret)(staffHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for customer accessing staff route, got %d", rec.Code)
		}
		if called {
			t.Error("handler should not have run")
		}
	})

	t.Run("RequireRole allows staff on staff endpoint", func(t *testing.T) {
		called := false
		staffHandler := middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/prescriptions/queue", nil)
		req.Header.Set("Authorization", "Bearer "+staffToken) // staff token
		rec := httptest.NewRecorder()

		middleware.AuthMiddleware(secret)(staffHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK for staff accessing staff route, got %d", rec.Code)
		}
		if !called {
			t.Error("staff handler should have run")
		}
	})
}
