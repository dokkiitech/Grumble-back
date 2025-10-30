package vibe

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// Repository defines persistence operations for vibes.
type Repository interface {
	// Create persists a new vibe and returns resulting counters.
	Create(ctx context.Context, vibe *Vibe) (*CreateResult, error)

	// Exists checks whether the user has already vibed the grumble.
	Exists(ctx context.Context, grumbleID shared.GrumbleID, userID shared.UserID) (bool, error)

	// CountByGrumble returns the number of vibes for a grumble.
	CountByGrumble(ctx context.Context, grumbleID shared.GrumbleID) (int, error)

	// FindByUser lists vibes created by a user, ordered by newest first.
	FindByUser(ctx context.Context, userID shared.UserID, limit int, offset int) ([]*Vibe, error)
}
