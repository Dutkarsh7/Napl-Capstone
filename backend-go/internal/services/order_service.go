package services

import (
	"smart-pharmacy/backend-go/internal/models"
)

// Re-export IsValidTransition for services package callers
func IsValidTransition(current, next models.OrderStatus) bool {
	return models.IsValidTransition(current, next)
}
