package grumble

import (
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// Grumble represents a user's complaint post (愚痴投稿)
type Grumble struct {
	GrumbleID      shared.GrumbleID
	UserID         shared.UserID
	Content        string
	ToxicLevel     shared.ToxicLevel
	VibeCount      int
	IsPurified     bool
	PostedAt       time.Time
	ExpiresAt      time.Time
	IsEventGrumble bool
}

// Validate checks if the grumble meets business rules
func (g *Grumble) Validate() error {
	// Content length: 1-280 characters
	if len(g.Content) == 0 {
		return &shared.ValidationError{
			Field:   "content",
			Message: "content cannot be empty",
		}
	}
	if len(g.Content) > 280 {
		return &shared.ValidationError{
			Field:   "content",
			Message: "content must be 280 characters or less",
		}
	}

	// Toxic level validation (1-5)
	if err := g.ToxicLevel.Validate(); err != nil {
		return err
	}

	// ExpiresAt must be after PostedAt
	if !g.ExpiresAt.After(g.PostedAt) {
		return &shared.ValidationError{
			Field:   "expires_at",
			Message: "expires_at must be after posted_at",
		}
	}

	return nil
}

// IsExpired checks if the grumble has passed its expiration time
func (g *Grumble) IsExpired() bool {
	return time.Now().After(g.ExpiresAt)
}

// CalculateTimeRemaining returns the duration until expiration
// Returns 0 if already expired
func (g *Grumble) CalculateTimeRemaining() time.Duration {
	remaining := time.Until(g.ExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// CanBePurified checks if the grumble has enough vibes for purification
func (g *Grumble) CanBePurified(threshold int) bool {
	return !g.IsPurified && g.VibeCount >= threshold
}

// Purify marks the grumble as purified (成仏)
func (g *Grumble) Purify() {
	g.IsPurified = true
}
