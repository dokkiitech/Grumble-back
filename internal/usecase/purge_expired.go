package usecase

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/logging"
)

// PurgeExpiredUseCase handles deletion of expired grumbles
type PurgeExpiredUseCase struct {
	grumbleRepo grumble.Repository
	logger      logging.Logger
}

// NewPurgeExpiredUseCase creates a new PurgeExpiredUseCase
func NewPurgeExpiredUseCase(grumbleRepo grumble.Repository, logger logging.Logger) *PurgeExpiredUseCase {
	return &PurgeExpiredUseCase{
		grumbleRepo: grumbleRepo,
		logger:      logger,
	}
}

// Purge deletes all expired grumbles (past 24 hours)
func (uc *PurgeExpiredUseCase) Purge(ctx context.Context) (int, error) {
	count, err := uc.grumbleRepo.DeleteExpired(ctx)
	if err != nil {
		uc.logger.ErrorContext(ctx, "Failed to purge expired grumbles", "error", err)
		return 0, err
	}

	if count > 0 {
		uc.logger.InfoContext(ctx, "Purged expired grumbles", "count", count)
	}

	return count, nil
}
