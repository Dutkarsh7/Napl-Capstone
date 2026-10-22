package tests

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/handlers"
	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
	"smart-pharmacy/backend-go/internal/services"
)

func TestPrescriptionUploadValidation(t *testing.T) {
	tempUploadDir, err := os.MkdirTemp("", "rx_test_uploads_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempUploadDir)

	store := db.NewMemoryStore()
	wp := services.NewWorkerPool(store, nil, nil, 1, 10) // Mock worker pool without starting network clients
	handler := handlers.NewPrescriptionsHandler(store, tempUploadDir, wp)

	user := &models.User{
		UserID: "user-upload-test",
		Email:  "uploader@test.local",
		Role:   models.RoleCustomer,
	}

	sendUpload := func(filename string, content []byte) *httptest.ResponseRecorder {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", filename)
		_, _ = part.Write(content)
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/prescriptions", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		handler.Upload(rec, req)
		return rec
	}

	t.Run("Rejects plain text file disguised with .jpg extension", func(t *testing.T) {
		fakeJpg := []byte("This is plain text and not a genuine image.")
		rec := sendUpload("fake_rx.jpg", fakeJpg)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request for text content with .jpg extension, got %d", rec.Code)
		}
	})

	t.Run("Accepts genuine PNG file by content sniffing", func(t *testing.T) {
		// Valid 1x1 transparent PNG magic bytes
		pngBytes := []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
			0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
			0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
		}

		rec := sendUpload("prescription.png", pngBytes)
		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201 Created for genuine PNG, got %d. Body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("Accepts genuine PDF file", func(t *testing.T) {
		pdfBytes := []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\ntrailer\n<< /Root 1 0 R >>\n%%EOF")
		rec := sendUpload("doctor_prescription.pdf", pdfBytes)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201 Created for PDF, got %d. Body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("Rejects file exceeding 5MB size limit", func(t *testing.T) {
		oversized := make([]byte, 6*1024*1024)
		copy(oversized, []byte("%PDF-1.4\n"))

		rec := sendUpload("oversized.pdf", oversized)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request for file > 5MB, got %d", rec.Code)
		}
	})
}
