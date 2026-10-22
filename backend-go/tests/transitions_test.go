package tests

import (
	"testing"

	"smart-pharmacy/backend-go/internal/models"
)

func TestOrderStatusTransitions(t *testing.T) {
	tests := []struct {
		name     string
		current  models.OrderStatus
		next     models.OrderStatus
		expected bool
	}{
		// Valid linear progressions
		{"PLACED to UNDER_REVIEW", models.StatusPlaced, models.StatusUnderReview, true},
		{"UNDER_REVIEW to CONFIRMED", models.StatusUnderReview, models.StatusConfirmed, true},
		{"CONFIRMED to READY", models.StatusConfirmed, models.StatusReady, true},
		{"READY to COMPLETED", models.StatusReady, models.StatusCompleted, true},

		// Valid cancellations
		{"PLACED to CANCELLED", models.StatusPlaced, models.StatusCancelled, true},
		{"UNDER_REVIEW to CANCELLED", models.StatusUnderReview, models.StatusCancelled, true},
		{"CONFIRMED to CANCELLED", models.StatusConfirmed, models.StatusCancelled, true},
		{"READY to CANCELLED", models.StatusReady, models.StatusCancelled, true},

		// Invalid leaps
		{"PLACED to COMPLETED (skip)", models.StatusPlaced, models.StatusCompleted, false},
		{"PLACED to READY (skip)", models.StatusPlaced, models.StatusReady, false},
		{"UNDER_REVIEW to COMPLETED (skip)", models.StatusUnderReview, models.StatusCompleted, false},

		// Terminal states cannot transition
		{"COMPLETED to CANCELLED", models.StatusCompleted, models.StatusCancelled, false},
		{"COMPLETED to PLACED", models.StatusCompleted, models.StatusPlaced, false},
		{"CANCELLED to PLACED", models.StatusCancelled, models.StatusPlaced, false},
		{"CANCELLED to COMPLETED", models.StatusCancelled, models.StatusCompleted, false},

		// Same state transition is invalid
		{"PLACED to PLACED", models.StatusPlaced, models.StatusPlaced, false},
		{"CONFIRMED to CONFIRMED", models.StatusConfirmed, models.StatusConfirmed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := models.IsValidTransition(tt.current, tt.next)
			if actual != tt.expected {
				t.Errorf("IsValidTransition(%s, %s) = %v; want %v", tt.current, tt.next, actual, tt.expected)
			}
		})
	}
}
