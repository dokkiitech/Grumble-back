package job

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/robfig/cron/v3"
)

// CronScheduler manages scheduled batch jobs
type CronScheduler struct {
	cron            *cron.Cron
	purgeExpiredJob *PurgeExpiredJob
	logger          logging.Logger
}

// NewCronScheduler creates a new CronScheduler
func NewCronScheduler(
	purgeExpiredJob *PurgeExpiredJob,
	logger logging.Logger,
) *CronScheduler {
	return &CronScheduler{
		cron:            cron.New(),
		purgeExpiredJob: purgeExpiredJob,
		logger:          logger,
	}
}

// Start initializes and starts all scheduled jobs
func (s *CronScheduler) Start() error {
	ctx := context.Background()
	// Run purge expired grumbles job every 5 minutes
	_, err := s.cron.AddFunc("*/5 * * * *", s.purgeExpiredJob.Run)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to schedule purge expired job", "error", err)
		return err
	}

	s.logger.InfoContext(ctx, "Starting cron scheduler")
	s.cron.Start()

	return nil
}

// Stop gracefully stops the cron scheduler
func (s *CronScheduler) Stop() {
	s.logger.InfoContext(context.Background(), "Stopping cron scheduler")
	s.cron.Stop()
}
