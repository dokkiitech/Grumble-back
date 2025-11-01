package user

import (
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// AnonymousUser represents an anonymous user account
type AnonymousUser struct {
	UserID       shared.UserID
	VirtuePoints int
	CreatedAt    time.Time
	ProfileTitle *string // Optional title like "今週の菩薩"
}

// Validate checks if the user meets business rules
func (u *AnonymousUser) Validate() error {
	// UserID must not be empty
	if u.UserID == "" {
		return &shared.ValidationError{
			Field:   "user_id",
			Message: "user_id cannot be empty",
		}
	}

	// Virtue points cannot be negative
	if u.VirtuePoints < 0 {
		return &shared.ValidationError{
			Field:   "virtue_points",
			Message: "virtue_points cannot be negative",
		}
	}

	// Profile title length check (if set)
	if u.ProfileTitle != nil && len(*u.ProfileTitle) > 50 {
		return &shared.ValidationError{
			Field:   "profile_title",
			Message: "profile_title must be 50 characters or less",
		}
	}

	return nil
}

// IncrementVirtuePoints adds points to the user's virtue score
func (u *AnonymousUser) IncrementVirtuePoints(points int) {
	u.VirtuePoints += points
}

// UpdateTitle sets or updates the user's profile title
func (u *AnonymousUser) UpdateTitle(title string) error {
	if len(title) > 50 {
		return &shared.ValidationError{
			Field:   "profile_title",
			Message: "profile_title must be 50 characters or less",
		}
	}
	u.ProfileTitle = &title
	return nil
}
