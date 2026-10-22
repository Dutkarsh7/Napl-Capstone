package models

// AllowedTransitions defines valid next statuses for each order status
var AllowedTransitions = map[OrderStatus][]OrderStatus{
	StatusPlaced: {
		StatusUnderReview,
		StatusCancelled,
	},
	StatusUnderReview: {
		StatusConfirmed,
		StatusCancelled,
	},
	StatusConfirmed: {
		StatusReady,
		StatusCancelled,
	},
	StatusReady: {
		StatusCompleted,
		StatusCancelled,
	},
	StatusCompleted: {}, // terminal state
	StatusCancelled: {}, // terminal state
}

// IsValidTransition returns true if transitioning from current to next is permitted
func IsValidTransition(current, next OrderStatus) bool {
	if current == next {
		return false
	}
	allowed, exists := AllowedTransitions[current]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}
