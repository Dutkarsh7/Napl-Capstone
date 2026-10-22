package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"smart-pharmacy/backend-go/internal/models"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, connString string) (*PostgresStore, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	cfg.MaxConns = 25
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

// Users

func (s *PostgresStore) CreateUser(ctx context.Context, u *models.User) error {
	if u.UserID == "" {
		u.UserID = uuid.NewString()
	}
	query := `
		INSERT INTO users (user_id, name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	_, err := s.pool.Exec(ctx, query, u.UserID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (s *PostgresStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT user_id, name, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	var u models.User
	err := s.pool.QueryRow(ctx, query, email).Scan(
		&u.UserID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (s *PostgresStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT user_id, name, email, password_hash, role, created_at
		FROM users
		WHERE user_id = $1
	`
	var u models.User
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&u.UserID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (s *PostgresStore) ListUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT user_id, name, email, role, created_at
		FROM users
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *PostgresStore) UpdateUserRole(ctx context.Context, id string, role models.UserRole) error {
	query := `UPDATE users SET role = $1 WHERE user_id = $2`
	ct, err := s.pool.Exec(ctx, query, role, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Medicines

func (s *PostgresStore) ListMedicines(ctx context.Context, query string, page, limit int) ([]models.Medicine, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM medicines
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' OR generic_name ILIKE '%' || $1 || '%')
	`
	if err := s.pool.QueryRow(ctx, countQuery, query).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT medicine_id, name, generic_name, description, price, stock_qty, requires_rx, created_at
		FROM medicines
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' OR generic_name ILIKE '%' || $1 || '%')
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.pool.Query(ctx, dataQuery, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var meds []models.Medicine
	for rows.Next() {
		var m models.Medicine
		if err := rows.Scan(&m.MedicineID, &m.Name, &m.GenericName, &m.Description, &m.Price, &m.StockQty, &m.RequiresRx, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		meds = append(meds, m)
	}
	return meds, total, rows.Err()
}

func (s *PostgresStore) GetMedicineByID(ctx context.Context, id string) (*models.Medicine, error) {
	query := `
		SELECT medicine_id, name, generic_name, description, price, stock_qty, requires_rx, created_at
		FROM medicines
		WHERE medicine_id = $1
	`
	var m models.Medicine
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&m.MedicineID, &m.Name, &m.GenericName, &m.Description, &m.Price, &m.StockQty, &m.RequiresRx, &m.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &m, err
}

func (s *PostgresStore) CreateMedicine(ctx context.Context, m *models.Medicine) error {
	if m.MedicineID == "" {
		m.MedicineID = uuid.NewString()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	query := `
		INSERT INTO medicines (medicine_id, name, generic_name, description, price, stock_qty, requires_rx, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query, m.MedicineID, m.Name, m.GenericName, m.Description, m.Price, m.StockQty, m.RequiresRx, m.CreatedAt)
	return err
}

func (s *PostgresStore) UpdateMedicine(ctx context.Context, m *models.Medicine) error {
	query := `
		UPDATE medicines
		SET name = $1, generic_name = $2, description = $3, price = $4, stock_qty = $5, requires_rx = $6
		WHERE medicine_id = $7
	`
	ct, err := s.pool.Exec(ctx, query, m.Name, m.GenericName, m.Description, m.Price, m.StockQty, m.RequiresRx, m.MedicineID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStore) UpdateMedicineStockAndPrice(ctx context.Context, id string, stockQty *int, price *float64) (*models.Medicine, error) {
	m, err := s.GetMedicineByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if stockQty != nil {
		m.StockQty = *stockQty
	}
	if price != nil {
		m.Price = *price
	}
	if err := s.UpdateMedicine(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Prescriptions

func (s *PostgresStore) CreatePrescription(ctx context.Context, p *models.Prescription) error {
	if p.PrescriptionID == "" {
		p.PrescriptionID = uuid.NewString()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	query := `
		INSERT INTO prescriptions (prescription_id, user_id, file_path, extracted_text, validation_status, review_status, review_note, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query, p.PrescriptionID, p.UserID, p.FilePath, p.ExtractedText, p.ValidationStatus, p.ReviewStatus, p.ReviewNote, p.CreatedAt)
	return err
}

func (s *PostgresStore) GetPrescriptionByID(ctx context.Context, id string) (*models.Prescription, error) {
	query := `
		SELECT prescription_id, user_id, file_path, extracted_text, validation_status, review_status, reviewed_by, review_note, created_at
		FROM prescriptions
		WHERE prescription_id = $1
	`
	var p models.Prescription
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&p.PrescriptionID, &p.UserID, &p.FilePath, &p.ExtractedText, &p.ValidationStatus, &p.ReviewStatus, &p.ReviewedBy, &p.ReviewNote, &p.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &p, err
}

func (s *PostgresStore) ListPrescriptionsByUser(ctx context.Context, userID string) ([]models.Prescription, error) {
	query := `
		SELECT prescription_id, user_id, file_path, extracted_text, validation_status, review_status, reviewed_by, review_note, created_at
		FROM prescriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Prescription
	for rows.Next() {
		var p models.Prescription
		if err := rows.Scan(&p.PrescriptionID, &p.UserID, &p.FilePath, &p.ExtractedText, &p.ValidationStatus, &p.ReviewStatus, &p.ReviewedBy, &p.ReviewNote, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (s *PostgresStore) ListPendingPrescriptions(ctx context.Context) ([]models.Prescription, error) {
	query := `
		SELECT prescription_id, user_id, file_path, extracted_text, validation_status, review_status, reviewed_by, review_note, created_at
		FROM prescriptions
		WHERE review_status = 'PENDING'
		ORDER BY created_at ASC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Prescription
	for rows.Next() {
		var p models.Prescription
		if err := rows.Scan(&p.PrescriptionID, &p.UserID, &p.FilePath, &p.ExtractedText, &p.ValidationStatus, &p.ReviewStatus, &p.ReviewedBy, &p.ReviewNote, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (s *PostgresStore) UpdatePrescriptionValidation(ctx context.Context, id string, extractedText string, status models.ValidationStatus) error {
	query := `
		UPDATE prescriptions
		SET extracted_text = $1, validation_status = $2
		WHERE prescription_id = $3
	`
	ct, err := s.pool.Exec(ctx, query, extractedText, status, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStore) ReviewPrescription(ctx context.Context, id string, reviewerID string, status models.ReviewStatus, note string) error {
	query := `
		UPDATE prescriptions
		SET review_status = $1, reviewed_by = $2, review_note = $3
		WHERE prescription_id = $4
	`
	ct, err := s.pool.Exec(ctx, query, status, reviewerID, note, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Orders

func (s *PostgresStore) CreateOrderWithItems(ctx context.Context, userID string, prescriptionID *string, items []models.CreateOrderItemRequest) (*models.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Check if prescription required
	requiresRx := false
	for _, item := range items {
		var isRx bool
		err := tx.QueryRow(ctx, `SELECT requires_rx FROM medicines WHERE medicine_id = $1`, item.MedicineID).Scan(&isRx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("medicine %s not found", item.MedicineID)
			}
			return nil, err
		}
		if isRx {
			requiresRx = true
		}
	}

	if requiresRx {
		if prescriptionID == nil || *prescriptionID == "" {
			return nil, ErrRxPrescriptionReq
		}
		var rxUserID string
		var rxStatus string
		err := tx.QueryRow(ctx, `SELECT user_id, review_status FROM prescriptions WHERE prescription_id = $1`, *prescriptionID).Scan(&rxUserID, &rxStatus)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrPrescriptionNotOwn
			}
			return nil, err
		}
		if rxUserID != userID || rxStatus != string(models.ReviewApproved) {
			return nil, ErrRxPrescriptionReq
		}
	}

	orderID := uuid.NewString()
	now := time.Now().UTC()
	var total float64
	var orderItems []models.OrderItem

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than zero")
		}

		// SELECT ... FOR UPDATE to lock row and avoid overselling
		var name string
		var price float64
		var stockQty int
		var rxReq bool
		lockQuery := `
			SELECT name, price, stock_qty, requires_rx
			FROM medicines
			WHERE medicine_id = $1
			FOR UPDATE
		`
		if err := tx.QueryRow(ctx, lockQuery, item.MedicineID).Scan(&name, &price, &stockQty, &rxReq); err != nil {
			return nil, err
		}

		if stockQty < item.Quantity {
			return nil, fmt.Errorf("%w: medicine %s (available: %d, requested: %d)", ErrInsufficientStock, name, stockQty, item.Quantity)
		}

		// Decrement stock
		newStock := stockQty - item.Quantity
		_, err := tx.Exec(ctx, `UPDATE medicines SET stock_qty = $1 WHERE medicine_id = $2`, newStock, item.MedicineID)
		if err != nil {
			return nil, err
		}

		lineTotal := price * float64(item.Quantity)
		total += lineTotal

		orderItem := models.OrderItem{
			ItemID:       uuid.NewString(),
			OrderID:      orderID,
			MedicineID:   item.MedicineID,
			MedicineName: name,
			Quantity:     item.Quantity,
			UnitPrice:    price,
			RequiresRx:   rxReq,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Insert order
	insertOrderQuery := `
		INSERT INTO orders (order_id, user_id, prescription_id, status, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, insertOrderQuery, orderID, userID, prescriptionID, models.StatusPlaced, total, now, now)
	if err != nil {
		return nil, err
	}

	// Insert order items
	insertItemQuery := `
		INSERT INTO order_items (item_id, order_id, medicine_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, it := range orderItems {
		_, err = tx.Exec(ctx, insertItemQuery, it.ItemID, it.OrderID, it.MedicineID, it.Quantity, it.UnitPrice)
		if err != nil {
			return nil, err
		}
	}

	// Insert initial status history
	insertHistQuery := `
		INSERT INTO order_status_history (id, order_id, status, changed_by, changed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, insertHistQuery, uuid.NewString(), orderID, models.StatusPlaced, userID, now)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Order{
		OrderID:        orderID,
		UserID:         userID,
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
	}, nil
}

func (s *PostgresStore) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	query := `
		SELECT o.order_id, o.user_id, u.name, u.email, o.prescription_id, o.status, o.total, o.created_at, o.updated_at
		FROM orders o
		JOIN users u ON o.user_id = u.user_id
		WHERE o.order_id = $1
	`
	var ord models.Order
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&ord.OrderID, &ord.UserID, &ord.CustomerName, &ord.CustomerEmail, &ord.PrescriptionID, &ord.Status, &ord.Total, &ord.CreatedAt, &ord.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Fetch items
	itemQuery := `
		SELECT oi.item_id, oi.order_id, oi.medicine_id, m.name, oi.quantity, oi.unit_price, m.requires_rx
		FROM order_items oi
		JOIN medicines m ON oi.medicine_id = m.medicine_id
		WHERE oi.order_id = $1
	`
	iRows, err := s.pool.Query(ctx, itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer iRows.Close()

	for iRows.Next() {
		var it models.OrderItem
		if err := iRows.Scan(&it.ItemID, &it.OrderID, &it.MedicineID, &it.MedicineName, &it.Quantity, &it.UnitPrice, &it.RequiresRx); err != nil {
			return nil, err
		}
		ord.Items = append(ord.Items, it)
	}

	// Fetch history
	histQuery := `
		SELECT id, order_id, status, changed_by, changed_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY changed_at ASC
	`
	hRows, err := s.pool.Query(ctx, histQuery, id)
	if err != nil {
		return nil, err
	}
	defer hRows.Close()

	for hRows.Next() {
		var h models.OrderStatusHistory
		if err := hRows.Scan(&h.ID, &h.OrderID, &h.Status, &h.ChangedBy, &h.ChangedAt); err != nil {
			return nil, err
		}
		ord.History = append(ord.History, h)
	}

	return &ord, nil
}

func (s *PostgresStore) ListOrdersByUser(ctx context.Context, userID string) ([]models.Order, error) {
	query := `
		SELECT order_id
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	var orders []models.Order
	for _, id := range ids {
		ord, err := s.GetOrderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *ord)
	}
	return orders, nil
}

func (s *PostgresStore) ListAllOrders(ctx context.Context) ([]models.Order, error) {
	query := `SELECT order_id FROM orders ORDER BY created_at DESC`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	var orders []models.Order
	for _, id := range ids {
		ord, err := s.GetOrderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *ord)
	}
	return orders, nil
}

func (s *PostgresStore) UpdateOrderStatus(ctx context.Context, id string, newStatus models.OrderStatus, changedBy string) (*models.Order, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var currentStatus models.OrderStatus
	err = tx.QueryRow(ctx, `SELECT status FROM orders WHERE order_id = $1 FOR UPDATE`, id).Scan(&currentStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if !models.IsValidTransition(currentStatus, newStatus) {
		return nil, fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidTransition, currentStatus, newStatus)
	}

	now := time.Now().UTC()

	// If cancelled, restore stock for all items in order
	if newStatus == models.StatusCancelled {
		itemRows, err := tx.Query(ctx, `SELECT medicine_id, quantity FROM order_items WHERE order_id = $1`, id)
		if err != nil {
			return nil, err
		}
		type restoreItem struct {
			medID string
			qty   int
		}
		var restores []restoreItem
		for itemRows.Next() {
			var r restoreItem
			if err := itemRows.Scan(&r.medID, &r.qty); err != nil {
				itemRows.Close()
				return nil, err
			}
			restores = append(restores, r)
		}
		itemRows.Close()

		for _, r := range restores {
			_, err = tx.Exec(ctx, `UPDATE medicines SET stock_qty = stock_qty + $1 WHERE medicine_id = $2`, r.qty, r.medID)
			if err != nil {
				return nil, err
			}
		}
	}

	_, err = tx.Exec(ctx, `UPDATE orders SET status = $1, updated_at = $2 WHERE order_id = $3`, newStatus, now, id)
	if err != nil {
		return nil, err
	}

	histID := uuid.NewString()
	_, err = tx.Exec(ctx, `INSERT INTO order_status_history (id, order_id, status, changed_by, changed_at) VALUES ($1, $2, $3, $4, $5)`,
		histID, id, newStatus, changedBy, now)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, id)
}
