package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smart-pharmacy/backend-go/internal/config"
	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/fsclient"
	"smart-pharmacy/backend-go/internal/handlers"
	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
	"smart-pharmacy/backend-go/internal/ocrclient"
	"smart-pharmacy/backend-go/internal/services"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting Smart Pharmacy Backend on port %s", cfg.Port)

	// Initialize Database Store (PostgreSQL if available, fallback to pre-seeded memory store if offline)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store, err := db.InitDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("[WARNING] Unable to connect to PostgreSQL (%v). Falling back to in-memory store for local execution.", err)
		mem := db.NewMemoryStore()
		db.SeedMemoryStore(mem)
		store = mem
	}
	defer store.Close()

	// Initialize Clients
	ocrClient := ocrclient.New(cfg.OCRServiceURL)
	fsClient := fsclient.New(cfg.ValidationServiceURL)

	// Initialize Worker Pool for background OCR and F# validation
	workerPool := services.NewWorkerPool(store, ocrClient, fsClient, cfg.WorkerPoolSize, cfg.WorkerQueueCap)
	workerPool.Start()
	defer workerPool.Stop()

	// Initialize Handlers
	authH := handlers.NewAuthHandler(store, cfg.JWTSecret)
	medH := handlers.NewMedicinesHandler(store)
	rxH := handlers.NewPrescriptionsHandler(store, cfg.UploadDir, workerPool)
	ordH := handlers.NewOrdersHandler(store)
	adminH := handlers.NewAdminHandler(store)

	// Router using standard Go 1.22+ ServeMux method patterns
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		handlers.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"disclaimer": "This platform assists pharmacy operations and does not provide medical advice.",
		})
	})

	// Auth
	mux.HandleFunc("POST /api/auth/register", authH.Register)
	mux.HandleFunc("POST /api/auth/login", authH.Login)
	mux.HandleFunc("GET /api/me", middleware.RequireAuth(authH.Me))

	// Catalogue
	mux.HandleFunc("GET /api/medicines", medH.List)
	mux.HandleFunc("GET /api/medicines/{id}", medH.GetByID)

	// Inventory (staff / admin)
	mux.HandleFunc("POST /api/medicines", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(medH.Create))
	mux.HandleFunc("PUT /api/medicines/{id}", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(medH.Update))
	mux.HandleFunc("PATCH /api/medicines/{id}/stock", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(medH.UpdateStock))

	// Prescriptions
	mux.HandleFunc("POST /api/prescriptions", middleware.RequireAuth(rxH.Upload))
	mux.HandleFunc("GET /api/prescriptions/mine", middleware.RequireAuth(rxH.Mine))
	mux.HandleFunc("GET /api/prescriptions/queue", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(rxH.Queue))
	mux.HandleFunc("GET /api/prescriptions/{id}", middleware.RequireAuth(rxH.GetByID))
	mux.HandleFunc("GET /api/prescriptions/{id}/file", middleware.RequireAuth(rxH.ServeFile))
	mux.HandleFunc("POST /api/prescriptions/{id}/review", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(rxH.Review))

	// Orders
	mux.HandleFunc("POST /api/orders", middleware.RequireAuth(ordH.Create))
	mux.HandleFunc("GET /api/orders/mine", middleware.RequireAuth(ordH.Mine))
	mux.HandleFunc("GET /api/orders/{id}", middleware.RequireAuth(ordH.GetByID))
	mux.HandleFunc("GET /api/orders", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(ordH.ListAll))
	mux.HandleFunc("PATCH /api/orders/{id}/status", middleware.RequireRole(models.RoleStaff, models.RoleAdmin)(ordH.UpdateStatus))

	// Admin
	mux.HandleFunc("GET /api/users", middleware.RequireRole(models.RoleAdmin)(adminH.ListUsers))
	mux.HandleFunc("PATCH /api/users/{id}/role", middleware.RequireRole(models.RoleAdmin)(adminH.UpdateRole))

	// Global Middlewares: Logger -> CORS -> Auth Context -> ServeMux
	handler := middleware.LoggerMiddleware(
		middleware.CORSMiddleware(cfg.AllowedOrigins)(
			middleware.AuthMiddleware(cfg.JWTSecret)(mux),
		),
	)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server shutdown listener
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("🚀 Smart Pharmacy Go Backend listening at http://localhost:%s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server ListenAndServe error: %v", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting successfully.")
}
