package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
)

type AuthHandler struct {
	store     db.Store
	jwtSecret string
}

func NewAuthHandler(store db.Store, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		JSONError(w, http.StatusBadRequest, "Name, valid email, and password of at least 6 characters are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         models.RoleCustomer,
		CreatedAt:    time.Now().UTC(),
	}

	if err := h.store.CreateUser(r.Context(), &user); err != nil {
		if errors.Is(err, db.ErrConflict) {
			JSONError(w, http.StatusConflict, "A user with this email already exists")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to create user account")
		return
	}

	token, err := middleware.GenerateToken(&user, h.jwtSecret)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Account created, but failed to generate token")
		return
	}

	JSONResponse(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || req.Password == "" {
		JSONError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		JSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		JSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := middleware.GenerateToken(user, h.jwtSecret)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	JSONResponse(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.store.GetUserByID(r.Context(), authUser.UserID)
	if err != nil {
		JSONError(w, http.StatusNotFound, "User not found")
		return
	}

	JSONResponse(w, http.StatusOK, user)
}
