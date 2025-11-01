package vibe

import (
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// Vibe represents an empathy reaction given to a grumble.
type Vibe struct {
	VibeID    shared.VibeID
	GrumbleID shared.GrumbleID
	UserID    shared.UserID
	Type      shared.VibeType
	VotedAt   time.Time
}

// Validate ensures the vibe adheres to domain rules.
func (v *Vibe) Validate() error {
	if v.GrumbleID == "" {
		return &shared.ValidationError{Field: "grumble_id", Message: "grumble_id is required"}
	}
	if v.UserID == "" {
		return &shared.ValidationError{Field: "user_id", Message: "user_id is required"}
	}
	if err := v.Type.Validate(); err != nil {
		return err
	}
	return nil
}

// CreateResult captures side effects of persisting a vibe.
type CreateResult struct {
	Vibe         *Vibe
	VibeCount    int
	VirtuePoints int
}
