package user

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// Repository defines the interface for user persistence
type Repository interface {
	// Create stores a new anonymous user
	Create(ctx context.Context, user *AnonymousUser) error

	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id shared.UserID) (*AnonymousUser, error)

	// Update updates an existing user
	Update(ctx context.Context, user *AnonymousUser) error

	// FindTopByVirtuePoints retrieves top users by virtue points for rankings
	FindTopByVirtuePoints(ctx context.Context, limit int) ([]*AnonymousUser, error)

	// IncrementVirtuePoints atomically increments a user's virtue points
	IncrementVirtuePoints(ctx context.Context, id shared.UserID, points int) error
}
