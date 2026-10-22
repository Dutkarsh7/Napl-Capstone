package db

import (
	"context"
	"errors"
	"smart-pharmacy/backend-go/internal/models"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrConflict            = errors.New("record already exists")
	ErrInsufficientStock   = errors.New("insufficient stock for medicine")
	ErrRxPrescriptionReq   = errors.New("prescription required for one or more items: approved prescription owned by user is mandatory")
	ErrInvalidTransition  = errors.New("invalid order status transition")
	ErrPrescriptionNotOwn  = errors.New("prescription not found or not owned by user")
)

type Store interface {
	// Users
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUserRole(ctx context.Context, id string, role models.UserRole) error

	// Medicines
	ListMedicines(ctx context.Context, query string, page, limit int) ([]models.Medicine, int, error)
	GetMedicineByID(ctx context.Context, id string) (*models.Medicine, error)
	CreateMedicine(ctx context.Context, m *models.Medicine) error
	UpdateMedicine(ctx context.Context, m *models.Medicine) error
	UpdateMedicineStockAndPrice(ctx context.Context, id string, stockQty *int, price *float64) (*models.Medicine, error)

	// Prescriptions
	CreatePrescription(ctx context.Context, p *models.Prescription) error
	GetPrescriptionByID(ctx context.Context, id string) (*models.Prescription, error)
	ListPrescriptionsByUser(ctx context.Context, userID string) ([]models.Prescription, error)
	ListPendingPrescriptions(ctx context.Context) ([]models.Prescription, error)
	UpdatePrescriptionValidation(ctx context.Context, id string, extractedText string, status models.ValidationStatus) error
	ReviewPrescription(ctx context.Context, id string, reviewerID string, status models.ReviewStatus, note string) error

	// Orders
	CreateOrderWithItems(ctx context.Context, userID string, prescriptionID *string, items []models.CreateOrderItemRequest) (*models.Order, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	ListOrdersByUser(ctx context.Context, userID string) ([]models.Order, error)
	ListAllOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, newStatus models.OrderStatus, changedBy string) (*models.Order, error)

	// Close
	Close()
}
