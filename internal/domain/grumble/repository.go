package grumble

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// TimelineFilter represents filtering options for timeline queries
type TimelineFilter struct {
	ToxicLevelMin  *shared.ToxicLevel // Minimum toxic level (inclusive)
	ToxicLevelMax  *shared.ToxicLevel // Maximum toxic level (inclusive)
	IsPurified     *bool              // Restrict to a specific purification state when provided
	ExcludeExpired bool               // Exclude expired grumbles
	UserID         *shared.UserID     // Filter by author user ID
	ViewerUserID   *shared.UserID     // Authenticated viewer for vibe state
	Limit          int                // Number of results to return
	Offset         int                // Number of results to skip
}

// Repository defines the interface for grumble persistence
type Repository interface {
	// Create stores a new grumble
	Create(ctx context.Context, grumble *Grumble) error

	// FindByID retrieves a grumble by its ID
	FindByID(ctx context.Context, id shared.GrumbleID) (*Grumble, error)

	// FindTimeline retrieves grumbles for the timeline with filtering
	FindTimeline(ctx context.Context, filter TimelineFilter) ([]*Grumble, error)

	// CountTimeline returns the total count of grumbles matching the filter
	CountTimeline(ctx context.Context, filter TimelineFilter) (int, error)

	// Update updates an existing grumble
	Update(ctx context.Context, grumble *Grumble) error

	// ArchiveExpired moves expired grumbles to archive table and removes them from main table
	ArchiveExpired(ctx context.Context) (int, error)

	// FindPurificationCandidates finds grumbles that meet purification threshold
	// but are not yet purified
	FindPurificationCandidates(ctx context.Context, threshold int) ([]*Grumble, error)

	// IncrementVibeCount atomically increments the vibe count for a grumble
	IncrementVibeCount(ctx context.Context, id shared.GrumbleID) error

	// FindArchivedTimeline retrieves grumbles from archive table for a specific date
	FindArchivedTimeline(ctx context.Context, filter TimelineFilter, targetDate time.Time) ([]*Grumble, error)

	// CountArchivedTimeline returns the total count of archived grumbles for a specific date
	CountArchivedTimeline(ctx context.Context, filter TimelineFilter, targetDate time.Time) (int, error)
}
