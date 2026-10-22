package services

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"smart-pharmacy/backend-go/internal/db"
	"smart-pharmacy/backend-go/internal/fsclient"
	"smart-pharmacy/backend-go/internal/models"
	"smart-pharmacy/backend-go/internal/ocrclient"
)

type PrescriptionJob struct {
	PrescriptionID string
	FilePath       string
}

type WorkerPool struct {
	store      db.Store
	ocrClient  *ocrclient.Client
	fsClient   *fsclient.Client
	jobQueue   chan PrescriptionJob
	workerCount int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewWorkerPool(store db.Store, ocr *ocrclient.Client, fs *fsclient.Client, workerCount, queueCap int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		store:       store,
		ocrClient:   ocr,
		fsClient:    fs,
		jobQueue:    make(chan PrescriptionJob, queueCap),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	log.Printf("[WorkerPool] Started %d background workers for prescription OCR & validation", wp.workerCount)
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
	log.Println("[WorkerPool] All workers stopped cleanly")
}

func (wp *WorkerPool) Enqueue(job PrescriptionJob) bool {
	select {
	case wp.jobQueue <- job:
		return true
	default:
		log.Printf("[WorkerPool] Queue full, failed to enqueue prescription %s", job.PrescriptionID)
		return false
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			wp.processJob(job)
		}
	}
}

func (wp *WorkerPool) processJob(job PrescriptionJob) {
	log.Printf("[WorkerPool] Processing prescription %s from %s", job.PrescriptionID, job.FilePath)

	jobCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// 1. OCR Extraction
	extractedText, err := wp.ocrClient.ExtractText(jobCtx, job.FilePath)
	if err != nil {
		log.Printf("[WorkerPool] OCR service unreachable or failed for %s: %v -> Defaulting to MANUAL_REVIEW", job.PrescriptionID, err)
		_ = wp.store.UpdatePrescriptionValidation(
			context.Background(),
			job.PrescriptionID,
			"(OCR service unavailable. Requires manual pharmacist transcription)",
			models.ValidationManualReview,
		)
		return
	}

	// 2. Parse text with naive regex parser
	parsed := ParsePrescriptionText(extractedText)

	// 3. Validate with F# microservice
	valResp, err := wp.fsClient.Validate(jobCtx, parsed)
	if err != nil {
		log.Printf("[WorkerPool] F# validation service unreachable or failed for %s: %v -> Defaulting to MANUAL_REVIEW", job.PrescriptionID, err)
		_ = wp.store.UpdatePrescriptionValidation(
			context.Background(),
			job.PrescriptionID,
			extractedText,
			models.ValidationManualReview,
		)
		return
	}

	// 4. Map F# status to domain validation status
	var valStatus models.ValidationStatus
	switch strings.ToLower(valResp.Status) {
	case "valid":
		valStatus = models.ValidationValid
	case "incomplete":
		valStatus = models.ValidationIncomplete
	case "invalid":
		valStatus = models.ValidationInvalid
	default:
		valStatus = models.ValidationManualReview
	}

	// Note: ReviewStatus remains PENDING. Validation is strictly human-in-the-loop assistance.
	if err := wp.store.UpdatePrescriptionValidation(context.Background(), job.PrescriptionID, extractedText, valStatus); err != nil {
		log.Printf("[WorkerPool] Error saving validation result for %s: %v", job.PrescriptionID, err)
		return
	}

	log.Printf("[WorkerPool] Successfully processed prescription %s -> validation_status: %s (awaiting pharmacist review)", job.PrescriptionID, valStatus)
}
