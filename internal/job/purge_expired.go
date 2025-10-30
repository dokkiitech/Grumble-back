package job

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
)

// PurgeExpiredJob is a cron job that deletes expired grumbles every 5 minutes
type PurgeExpiredJob struct {
	purgeExpiredUC *usecase.PurgeExpiredUseCase
	logger         logging.Logger
}

// NewPurgeExpiredJob creates a new PurgeExpiredJob
func NewPurgeExpiredJob(
	purgeExpiredUC *usecase.PurgeExpiredUseCase,
	logger logging.Logger,
) *PurgeExpiredJob {
	return &PurgeExpiredJob{
		purgeExpiredUC: purgeExpiredUC,
		logger:         logger,
	}
}

// Run executes the purge expired grumbles job
func (j *PurgeExpiredJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	j.logger.InfoContext(ctx, "Starting purge expired grumbles job")

	count, err := j.purgeExpiredUC.Purge(ctx)
	if err != nil {
		j.logger.ErrorContext(ctx, "Purge expired job failed", "error", err)
		return
	}

	if count > 0 {
		j.logger.InfoContext(ctx, "Purge expired job completed", "deleted_count", count)
	}
}
