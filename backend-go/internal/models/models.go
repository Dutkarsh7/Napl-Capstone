package models

import (
	"time"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleStaff    UserRole = "staff"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Medicine struct {
	MedicineID  string    `json:"medicine_id"`
	Name        string    `json:"name"`
	GenericName string    `json:"generic_name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	StockQty    int       `json:"stock_qty"`
	RequiresRx  bool      `json:"requires_rx"`
	CreatedAt   time.Time `json:"created_at"`
}

type ValidationStatus string

const (
	ValidationPending      ValidationStatus = "PENDING"
	ValidationValid        ValidationStatus = "VALID"
	ValidationIncomplete   ValidationStatus = "INCOMPLETE"
	ValidationInvalid      ValidationStatus = "INVALID"
	ValidationManualReview ValidationStatus = "MANUAL_REVIEW"
)

type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "PENDING"
	ReviewApproved ReviewStatus = "APPROVED"
	ReviewRejected ReviewStatus = "REJECTED"
)

type Prescription struct {
	PrescriptionID   string           `json:"prescription_id"`
	UserID           string           `json:"user_id"`
	FilePath         string           `json:"file_path"`
	ExtractedText    string           `json:"extracted_text"`
	ValidationStatus ValidationStatus `json:"validation_status"`
	ReviewStatus     ReviewStatus     `json:"review_status"`
	ReviewedBy       *string          `json:"reviewed_by,omitempty"`
	ReviewNote       string           `json:"review_note"`
	CreatedAt        time.Time        `json:"created_at"`
}

type OrderStatus string

const (
	StatusPlaced      OrderStatus = "PLACED"
	StatusUnderReview OrderStatus = "UNDER_REVIEW"
	StatusConfirmed   OrderStatus = "CONFIRMED"
	StatusReady       OrderStatus = "READY"
	StatusCompleted   OrderStatus = "COMPLETED"
	StatusCancelled   OrderStatus = "CANCELLED"
)

type OrderItem struct {
	ItemID       string    `json:"item_id"`
	OrderID      string    `json:"order_id"`
	MedicineID   string    `json:"medicine_id"`
	MedicineName string    `json:"medicine_name,omitempty"`
	Quantity     int       `json:"quantity"`
	UnitPrice    float64   `json:"unit_price"`
	RequiresRx   bool      `json:"requires_rx,omitempty"`
}

type OrderStatusHistory struct {
	ID        string      `json:"id"`
	OrderID   string      `json:"order_id"`
	Status    OrderStatus `json:"status"`
	ChangedBy *string     `json:"changed_by,omitempty"`
	ChangedAt time.Time   `json:"changed_at"`
}

type Order struct {
	OrderID        string               `json:"order_id"`
	UserID         string               `json:"user_id"`
	CustomerName   string               `json:"customer_name,omitempty"`
	CustomerEmail  string               `json:"customer_email,omitempty"`
	PrescriptionID *string              `json:"prescription_id,omitempty"`
	Status         OrderStatus          `json:"status"`
	Total          float64              `json:"total"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Items          []OrderItem          `json:"items,omitempty"`
	History        []OrderStatusHistory `json:"history,omitempty"`
}

// Request and Response payloads

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateOrderItemRequest struct {
	MedicineID string `json:"medicine_id"`
	Quantity   int    `json:"quantity"`
}

type CreateOrderRequest struct {
	PrescriptionID *string                  `json:"prescription_id,omitempty"`
	Items          []CreateOrderItemRequest `json:"items"`
}

type ReviewPrescriptionRequest struct {
	Action string `json:"action"` // "APPROVE" or "REJECT"
	Note   string `json:"note"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status"`
}

type UpdateStockRequest struct {
	StockQty *int     `json:"stock_qty,omitempty"`
	Price    *float64 `json:"price,omitempty"`
}

type UpdateUserRoleRequest struct {
	Role UserRole `json:"role"`
}
