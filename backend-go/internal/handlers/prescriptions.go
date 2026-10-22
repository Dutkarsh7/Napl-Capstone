package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/middleware"
	"smart-pharmacy/backend-go/internal/models"
	"smart-pharmacy/backend-go/internal/services"
)

type PrescriptionsHandler struct {
	store      db.Store
	uploadDir  string
	workerPool *services.WorkerPool
}

func NewPrescriptionsHandler(store db.Store, uploadDir string, wp *services.WorkerPool) *PrescriptionsHandler {
	return &PrescriptionsHandler{
		store:      store,
		uploadDir:  uploadDir,
		workerPool: wp,
	}
}

// 5 MB max file size
const maxUploadSize = 5 * 1024 * 1024

func (h *PrescriptionsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+1024*1024)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		JSONError(w, http.StatusBadRequest, "File size exceeds 5MB limit or invalid form")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Prescription file is required (form field 'file')")
		return
	}
	defer file.Close()

	if handler.Size > maxUploadSize {
		JSONError(w, http.StatusBadRequest, "File size exceeds 5 MB limit")
		return
	}

	// Sniff first 512 bytes to verify real MIME type
	sniffBuffer := make([]byte, 512)
	n, err := file.Read(sniffBuffer)
	if err != nil && err != io.EOF {
		JSONError(w, http.StatusInternalServerError, "Failed to read uploaded file")
		return
	}
	sniffBuffer = sniffBuffer[:n]

	// Reset file pointer to beginning
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			JSONError(w, http.StatusInternalServerError, "Failed to process uploaded file")
			return
		}
	} else {
		JSONError(w, http.StatusInternalServerError, "Stream processing error")
		return
	}

	contentType := http.DetectContentType(sniffBuffer)
	ext := strings.ToLower(filepath.Ext(handler.Filename))

	var validExt string
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		validExt = ".jpg"
	case strings.HasPrefix(contentType, "image/png"):
		validExt = ".png"
	case strings.HasPrefix(contentType, "application/pdf") || (ext == ".pdf" && strings.Contains(string(sniffBuffer), "%PDF")):
		validExt = ".pdf"
	default:
		JSONError(w, http.StatusBadRequest, fmt.Sprintf("Unsupported file type: detected %s. Only JPEG, PNG, and PDF files are allowed", contentType))
		return
	}

	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to prepare storage directory")
		return
	}

	randomFilename := fmt.Sprintf("%s%s", uuid.NewString(), validExt)
	destPath := filepath.Join(h.uploadDir, randomFilename)

	dst, err := os.Create(destPath)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to save prescription file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to write file to disk")
		return
	}

	prescriptionID := uuid.NewString()
	now := time.Now().UTC()

	prescription := models.Prescription{
		PrescriptionID:   prescriptionID,
		UserID:           authUser.UserID,
		FilePath:         destPath,
		ExtractedText:    "",
		ValidationStatus: models.ValidationPending,
		ReviewStatus:     models.ReviewPending,
		ReviewNote:       "",
		CreatedAt:        now,
	}

	if err := h.store.CreatePrescription(r.Context(), &prescription); err != nil {
		_ = os.Remove(destPath)
		JSONError(w, http.StatusInternalServerError, "Failed to record prescription")
		return
	}

	// Dispatch asynchronous OCR & F# validation job to worker pool
	h.workerPool.Enqueue(services.PrescriptionJob{
		PrescriptionID: prescriptionID,
		FilePath:       destPath,
	})

	JSONResponse(w, http.StatusCreated, prescription)
}

func (h *PrescriptionsHandler) Mine(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	prescriptions, err := h.store.ListPrescriptionsByUser(r.Context(), authUser.UserID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve prescriptions")
		return
	}

	JSONResponse(w, http.StatusOK, prescriptions)
}

func (h *PrescriptionsHandler) Queue(w http.ResponseWriter, r *http.Request) {
	prescriptions, err := h.store.ListPendingPrescriptions(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to load prescription review queue")
		return
	}

	JSONResponse(w, http.StatusOK, prescriptions)
}

func (h *PrescriptionsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Prescription ID is required")
		return
	}

	rx, err := h.store.GetPrescriptionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Prescription not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve prescription")
		return
	}

	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Authorization check: owner, staff, or admin
	if rx.UserID != authUser.UserID && authUser.Role != models.RoleStaff && authUser.Role != models.RoleAdmin {
		JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	JSONResponse(w, http.StatusOK, rx)
}

func (h *PrescriptionsHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Prescription ID is required")
		return
	}

	rx, err := h.store.GetPrescriptionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Prescription not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve prescription")
		return
	}

	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if rx.UserID != authUser.UserID && authUser.Role != models.RoleStaff && authUser.Role != models.RoleAdmin {
		JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	http.ServeFile(w, r, rx.FilePath)
}

func (h *PrescriptionsHandler) Review(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Prescription ID is required")
		return
	}

	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok || authUser == nil {
		JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	var req models.ReviewPrescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var status models.ReviewStatus
	switch strings.ToUpper(strings.TrimSpace(req.Action)) {
	case "APPROVE":
		status = models.ReviewApproved
	case "REJECT":
		status = models.ReviewRejected
	default:
		JSONError(w, http.StatusBadRequest, "Action must be APPROVE or REJECT")
		return
	}

	if err := h.store.ReviewPrescription(r.Context(), id, authUser.UserID, status, strings.TrimSpace(req.Note)); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			JSONError(w, http.StatusNotFound, "Prescription not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "Failed to submit review")
		return
	}

	rx, err := h.store.GetPrescriptionByID(r.Context(), id)
	if err != nil {
		JSONResponse(w, http.StatusOK, map[string]string{"message": "Prescription reviewed successfully"})
		return
	}

	JSONResponse(w, http.StatusOK, rx)
}
