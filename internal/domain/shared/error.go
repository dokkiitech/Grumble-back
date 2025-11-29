package shared

import "fmt"

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Entity string
	ID     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Entity, e.ID)
}

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// DuplicateVibeError represents an attempt to give a duplicate vibe
type DuplicateVibeError struct {
	GrumbleID string
	UserID    string
}

func (e *DuplicateVibeError) Error() string {
	return fmt.Sprintf("user %s already gave vibe to grumble %s", e.UserID, e.GrumbleID)
}

// UnauthorizedError represents an authentication failure
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

// InternalError represents an internal system error
type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("internal error: %s", e.Message)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

// InappropriateContentError represents content that violates moderation rules
type InappropriateContentError struct {
	Reason string
}

func (e *InappropriateContentError) Error() string {
	return e.Reason
}
