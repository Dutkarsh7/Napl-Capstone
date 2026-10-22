package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"smart-pharmacy/backend-go/internal/models"
)

type MemoryStore struct {
	mu            sync.RWMutex
	users         map[string]*models.User
	usersByEmail  map[string]string
	medicines     map[string]*models.Medicine
	prescriptions map[string]*models.Prescription
	orders        map[string]*models.Order
	orderItems    map[string][]models.OrderItem
	orderHistory  map[string][]models.OrderStatusHistory
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:         make(map[string]*models.User),
		usersByEmail:  make(map[string]string),
		medicines:     make(map[string]*models.Medicine),
		prescriptions: make(map[string]*models.Prescription),
		orders:        make(map[string]*models.Order),
		orderItems:    make(map[string][]models.OrderItem),
		orderHistory:  make(map[string][]models.OrderStatusHistory),
	}
}

func (m *MemoryStore) Close() {}

// Users

func (m *MemoryStore) CreateUser(ctx context.Context, u *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	email := strings.ToLower(strings.TrimSpace(u.Email))
	if _, exists := m.usersByEmail[email]; exists {
		return ErrConflict
	}
	if u.UserID == "" {
		u.UserID = uuid.NewString()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	cp := *u
	cp.Email = email
	m.users[u.UserID] = &cp
	m.usersByEmail[email] = u.UserID
	return nil
}

func (m *MemoryStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.usersByEmail[strings.ToLower(strings.TrimSpace(email))]
	if !exists {
		return nil, ErrNotFound
	}
	u := *m.users[id]
	return &u, nil
}

func (m *MemoryStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, exists := m.users[id]
	if !exists {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *MemoryStore) ListUsers(ctx context.Context) ([]models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []models.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *MemoryStore) UpdateUserRole(ctx context.Context, id string, role models.UserRole) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, exists := m.users[id]
	if !exists {
		return ErrNotFound
	}
	u.Role = role
	return nil
}

// Medicines

func (m *MemoryStore) ListMedicines(ctx context.Context, query string, page, limit int) ([]models.Medicine, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []models.Medicine
	q := strings.ToLower(strings.TrimSpace(query))
	for _, med := range m.medicines {
		if q == "" || strings.Contains(strings.ToLower(med.Name), q) || strings.Contains(strings.ToLower(med.GenericName), q) {
			filtered = append(filtered, *med)
		}
	}

	total := len(filtered)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	start := (page - 1) * limit
	if start >= total {
		return []models.Medicine{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

func (m *MemoryStore) GetMedicineByID(ctx context.Context, id string) (*models.Medicine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	med, exists := m.medicines[id]
	if !exists {
		return nil, ErrNotFound
	}
	cp := *med
	return &cp, nil
}

func (m *MemoryStore) CreateMedicine(ctx context.Context, med *models.Medicine) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if med.MedicineID == "" {
		med.MedicineID = uuid.NewString()
	}
	if med.CreatedAt.IsZero() {
		med.CreatedAt = time.Now().UTC()
	}
	cp := *med
	m.medicines[med.MedicineID] = &cp
	return nil
}

func (m *MemoryStore) UpdateMedicine(ctx context.Context, med *models.Medicine) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.medicines[med.MedicineID]; !exists {
		return ErrNotFound
	}
	cp := *med
	m.medicines[med.MedicineID] = &cp
	return nil
}

func (m *MemoryStore) UpdateMedicineStockAndPrice(ctx context.Context, id string, stockQty *int, price *float64) (*models.Medicine, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	med, exists := m.medicines[id]
	if !exists {
		return nil, ErrNotFound
	}
	if stockQty != nil {
		med.StockQty = *stockQty
	}
	if price != nil {
		med.Price = *price
	}
	cp := *med
	return &cp, nil
}

// Prescriptions

func (m *MemoryStore) CreatePrescription(ctx context.Context, p *models.Prescription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p.PrescriptionID == "" {
		p.PrescriptionID = uuid.NewString()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	cp := *p
	m.prescriptions[p.PrescriptionID] = &cp
	return nil
}

func (m *MemoryStore) GetPrescriptionByID(ctx context.Context, id string) (*models.Prescription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.prescriptions[id]
	if !exists {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *MemoryStore) ListPrescriptionsByUser(ctx context.Context, userID string) ([]models.Prescription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []models.Prescription
	for _, p := range m.prescriptions {
		if p.UserID == userID {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *MemoryStore) ListPendingPrescriptions(ctx context.Context) ([]models.Prescription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []models.Prescription
	for _, p := range m.prescriptions {
		if p.ReviewStatus == models.ReviewPending {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *MemoryStore) UpdatePrescriptionValidation(ctx context.Context, id string, extractedText string, status models.ValidationStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.prescriptions[id]
	if !exists {
		return ErrNotFound
	}
	p.ExtractedText = extractedText
	p.ValidationStatus = status
	return nil
}

func (m *MemoryStore) ReviewPrescription(ctx context.Context, id string, reviewerID string, status models.ReviewStatus, note string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.prescriptions[id]
	if !exists {
		return ErrNotFound
	}
	p.ReviewStatus = status
	p.ReviewedBy = &reviewerID
	p.ReviewNote = note
	return nil
}

// Orders

func (m *MemoryStore) CreateOrderWithItems(ctx context.Context, userID string, prescriptionID *string, items []models.CreateOrderItemRequest) (*models.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// 1. Verify Rx items
	requiresRx := false
	for _, it := range items {
		med, exists := m.medicines[it.MedicineID]
		if !exists {
			return nil, fmt.Errorf("medicine %s not found", it.MedicineID)
		}
		if med.RequiresRx {
			requiresRx = true
		}
	}

	if requiresRx {
		if prescriptionID == nil || *prescriptionID == "" {
			return nil, ErrRxPrescriptionReq
		}
		rx, exists := m.prescriptions[*prescriptionID]
		if !exists || rx.UserID != userID {
			return nil, ErrPrescriptionNotOwn
		}
		if rx.ReviewStatus != models.ReviewApproved {
			return nil, ErrRxPrescriptionReq
		}
	}

	// 2. Verify stock
	for _, it := range items {
		if it.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than zero")
		}
		med := m.medicines[it.MedicineID]
		if med.StockQty < it.Quantity {
			return nil, fmt.Errorf("%w: medicine %s (available: %d, requested: %d)", ErrInsufficientStock, med.Name, med.StockQty, it.Quantity)
		}
	}

	// 3. Deduct stock & calculate total
	orderID := uuid.NewString()
	now := time.Now().UTC()
	var total float64
	var orderItems []models.OrderItem

	for _, it := range items {
		med := m.medicines[it.MedicineID]
		med.StockQty -= it.Quantity
		lineTotal := med.Price * float64(it.Quantity)
		total += lineTotal

		orderItems = append(orderItems, models.OrderItem{
			ItemID:       uuid.NewString(),
			OrderID:      orderID,
			MedicineID:   it.MedicineID,
			MedicineName: med.Name,
			Quantity:     it.Quantity,
			UnitPrice:    med.Price,
			RequiresRx:   med.RequiresRx,
		})
	}

	userObj := m.users[userID]
	customerName := ""
	customerEmail := ""
	if userObj != nil {
		customerName = userObj.Name
		customerEmail = userObj.Email
	}

	ord := &models.Order{
		OrderID:        orderID,
		UserID:         userID,
		CustomerName:   customerName,
		CustomerEmail:  customerEmail,
		PrescriptionID: prescriptionID,
		Status:         models.StatusPlaced,
		Total:          total,
		CreatedAt:      now,
		UpdatedAt:      now,
		Items:          orderItems,
		History: []models.OrderStatusHistory{
			{
				ID:        uuid.NewString(),
				OrderID:   orderID,
				Status:    models.StatusPlaced,
				ChangedBy: &userID,
				ChangedAt: now,
			},
		},
	}

	m.orders[orderID] = ord
	m.orderItems[orderID] = orderItems
	m.orderHistory[orderID] = ord.History

	cp := *ord
	return &cp, nil
}

func (m *MemoryStore) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ord, exists := m.orders[id]
	if !exists {
		return nil, ErrNotFound
	}
	cp := *ord
	cp.Items = m.orderItems[id]
	cp.History = m.orderHistory[id]
	return &cp, nil
}

func (m *MemoryStore) ListOrdersByUser(ctx context.Context, userID string) ([]models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []models.Order
	for _, ord := range m.orders {
		if ord.UserID == userID {
			cp := *ord
			cp.Items = m.orderItems[ord.OrderID]
			cp.History = m.orderHistory[ord.OrderID]
			list = append(list, cp)
		}
	}
	return list, nil
}

func (m *MemoryStore) ListAllOrders(ctx context.Context) ([]models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []models.Order
	for _, ord := range m.orders {
		cp := *ord
		cp.Items = m.orderItems[ord.OrderID]
		cp.History = m.orderHistory[ord.OrderID]
		list = append(list, cp)
	}
	return list, nil
}

func (m *MemoryStore) UpdateOrderStatus(ctx context.Context, id string, newStatus models.OrderStatus, changedBy string) (*models.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ord, exists := m.orders[id]
	if !exists {
		return nil, ErrNotFound
	}

	if !models.IsValidTransition(ord.Status, newStatus) {
		return nil, fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidTransition, ord.Status, newStatus)
	}

	now := time.Now().UTC()

	// If cancelled, restore stock
	if newStatus == models.StatusCancelled {
		for _, it := range m.orderItems[id] {
			if med, ok := m.medicines[it.MedicineID]; ok {
				med.StockQty += it.Quantity
			}
		}
	}

	ord.Status = newStatus
	ord.UpdatedAt = now

	hist := models.OrderStatusHistory{
		ID:        uuid.NewString(),
		OrderID:   id,
		Status:    newStatus,
		ChangedBy: &changedBy,
		ChangedAt: now,
	}
	m.orderHistory[id] = append(m.orderHistory[id], hist)
	ord.History = m.orderHistory[id]

	cp := *ord
	cp.Items = m.orderItems[id]
	return &cp, nil
}
