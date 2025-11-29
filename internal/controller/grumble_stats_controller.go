package controller

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
)

// GrumbleStatsController handles stats-related requests.
type GrumbleStatsController struct {
	uc     *usecase.GrumbleStatsUseCase
	logger logging.Logger
}

// NewGrumbleStatsController creates a new controller.
func NewGrumbleStatsController(uc *usecase.GrumbleStatsUseCase, logger logging.Logger) *GrumbleStatsController {
	return &GrumbleStatsController{uc: uc, logger: logger}
}

// StatsInput represents params for total stats.
type StatsInput struct {
	Granularity grumble.Granularity
	From        *time.Time
	To          *time.Time
	TZ          string
}

// StatsByToxicInput represents params for toxic-level stats.
type StatsByToxicInput struct {
	StatsInput
	ToxicLevel *int
}

// GetStats returns aggregated stats.
func (c *GrumbleStatsController) GetStats(ctx context.Context, in StatsInput) ([]grumble.StatsRow, error) {
	req := usecase.StatsRequest{
		Granularity: in.Granularity,
		From:        in.From,
		To:          in.To,
		TZ:          in.TZ,
	}
	return c.uc.Get(ctx, req)
}

// GetStatsByToxic returns aggregated stats grouped by toxic level.
func (c *GrumbleStatsController) GetStatsByToxic(ctx context.Context, in StatsByToxicInput) ([]grumble.StatsRow, error) {
	req := usecase.StatsByToxicRequest{
		StatsRequest: usecase.StatsRequest{
			Granularity: in.Granularity,
			From:        in.From,
			To:          in.To,
			TZ:          in.TZ,
		},
		ToxicLevel: in.ToxicLevel,
	}
	return c.uc.GetByToxic(ctx, req)
}
