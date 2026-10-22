package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/models"
)

func TestOrderStockAndRxLogic(t *testing.T) {
	ctx := context.Background()

	setupStore := func() (*db.MemoryStore, string, string, string) {
		store := db.NewMemoryStore()

		// Seed a customer
		cust := models.User{
			UserID:       "cust-123",
			Name:         "Alice Tester",
			Email:        "alice@test.local",
			PasswordHash: "dummyhash",
			Role:         models.RoleCustomer,
			CreatedAt:    time.Now().UTC(),
		}
		_ = store.CreateUser(ctx, &cust)

		// Seed OTC medicine: Stock 10
		otcMed := models.Medicine{
			MedicineID:  "otc-1",
			Name:        "Paracetamol 500mg",
			GenericName: "Paracetamol",
			Price:       5.00,
			StockQty:    10,
			RequiresRx:  false,
		}
		_ = store.CreateMedicine(ctx, &otcMed)

		// Seed Rx medicine: Stock 10
		rxMed := models.Medicine{
			MedicineID:  "rx-1",
			Name:        "Amoxicillin 500mg",
			GenericName: "Amoxicillin",
			Price:       15.00,
			StockQty:    10,
			RequiresRx:  true,
		}
		_ = store.CreateMedicine(ctx, &rxMed)

		return store, cust.UserID, otcMed.MedicineID, rxMed.MedicineID
	}

	t.Run("Successful OTC order decrements stock accurately", func(t *testing.T) {
		store, userID, otcID, _ := setupStore()

		items := []models.CreateOrderItemRequest{
			{MedicineID: otcID, Quantity: 3},
		}

		order, err := store.CreateOrderWithItems(ctx, userID, nil, items)
		if err != nil {
			t.Fatalf("expected successful order, got error: %v", err)
		}

		if order.Total != 15.00 {
			t.Errorf("expected total 15.00, got %.2f", order.Total)
		}
		if order.Status != models.StatusPlaced {
			t.Errorf("expected status PLACED, got %s", order.Status)
		}

		// Verify stock decremented: 10 - 3 = 7
		med, _ := store.GetMedicineByID(ctx, otcID)
		if med.StockQty != 7 {
			t.Errorf("expected remaining stock 7, got %d", med.StockQty)
		}
	})

	t.Run("Order exceeds stock fails with ErrInsufficientStock without modifying stock", func(t *testing.T) {
		store, userID, otcID, _ := setupStore()

		items := []models.CreateOrderItemRequest{
			{MedicineID: otcID, Quantity: 20}, // Stock is only 10
		}

		_, err := store.CreateOrderWithItems(ctx, userID, nil, items)
		if !errors.Is(err, db.ErrInsufficientStock) {
			t.Fatalf("expected ErrInsufficientStock, got: %v", err)
		}

		// Verify stock unchanged
		med, _ := store.GetMedicineByID(ctx, otcID)
		if med.StockQty != 10 {
			t.Errorf("expected stock to remain 10, got %d", med.StockQty)
		}
	})

	t.Run("Ordering Rx medicine without prescription fails with ErrRxPrescriptionReq", func(t *testing.T) {
		store, userID, _, rxID := setupStore()

		items := []models.CreateOrderItemRequest{
			{MedicineID: rxID, Quantity: 2},
		}

		_, err := store.CreateOrderWithItems(ctx, userID, nil, items)
		if !errors.Is(err, db.ErrRxPrescriptionReq) {
			t.Fatalf("expected ErrRxPrescriptionReq, got: %v", err)
		}
	})

	t.Run("Ordering Rx medicine with pending prescription fails", func(t *testing.T) {
		store, userID, _, rxID := setupStore()

		// Create prescription with PENDING review status
		rx := models.Prescription{
			PrescriptionID:   "rx-doc-1",
			UserID:           userID,
			FilePath:         "/uploads/test.png",
			ValidationStatus: models.ValidationValid,
			ReviewStatus:     models.ReviewPending, // Not yet approved
			CreatedAt:        time.Now().UTC(),
		}
		_ = store.CreatePrescription(ctx, &rx)

		items := []models.CreateOrderItemRequest{
			{MedicineID: rxID, Quantity: 2},
		}

		rxIDStr := rx.PrescriptionID
		_, err := store.CreateOrderWithItems(ctx, userID, &rxIDStr, items)
		if !errors.Is(err, db.ErrRxPrescriptionReq) {
			t.Fatalf("expected ErrRxPrescriptionReq for unapproved rx, got: %v", err)
		}
	})

	t.Run("Ordering Rx medicine with approved prescription succeeds and decrements stock", func(t *testing.T) {
		store, userID, _, rxID := setupStore()

		// Create prescription with APPROVED status
		rx := models.Prescription{
			PrescriptionID:   "rx-approved-1",
			UserID:           userID,
			FilePath:         "/uploads/test.png",
			ValidationStatus: models.ValidationValid,
			ReviewStatus:     models.ReviewApproved,
			CreatedAt:        time.Now().UTC(),
		}
		_ = store.CreatePrescription(ctx, &rx)

		items := []models.CreateOrderItemRequest{
			{MedicineID: rxID, Quantity: 4},
		}

		rxIDStr := rx.PrescriptionID
		order, err := store.CreateOrderWithItems(ctx, userID, &rxIDStr, items)
		if err != nil {
			t.Fatalf("expected success with approved rx, got: %v", err)
		}

		if order.Total != 60.00 {
			t.Errorf("expected total 60.00, got %.2f", order.Total)
		}

		med, _ := store.GetMedicineByID(ctx, rxID)
		if med.StockQty != 6 {
			t.Errorf("expected remaining stock 6, got %d", med.StockQty)
		}
	})

	t.Run("Cancelling order restores stock", func(t *testing.T) {
		store, userID, otcID, _ := setupStore()

		items := []models.CreateOrderItemRequest{
			{MedicineID: otcID, Quantity: 5},
		}

		order, err := store.CreateOrderWithItems(ctx, userID, nil, items)
		if err != nil {
			t.Fatalf("setup order failed: %v", err)
		}

		med, _ := store.GetMedicineByID(ctx, otcID)
		if med.StockQty != 5 {
			t.Fatalf("expected stock 5 after order, got %d", med.StockQty)
		}

		// Cancel the order
		updated, err := store.UpdateOrderStatus(ctx, order.OrderID, models.StatusCancelled, userID)
		if err != nil {
			t.Fatalf("expected successful cancellation, got: %v", err)
		}
		if updated.Status != models.StatusCancelled {
			t.Errorf("expected status CANCELLED, got %s", updated.Status)
		}

		// Verify stock was restored to 10
		medRestored, _ := store.GetMedicineByID(ctx, otcID)
		if medRestored.StockQty != 10 {
			t.Errorf("expected restored stock 10, got %d", medRestored.StockQty)
		}
	})
}
