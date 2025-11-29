package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
)

// GrumbleStatsUseCase handles stats retrieval with default period computation.
type GrumbleStatsUseCase struct {
	repo               grumble.Repository
	defaultTZ          string
	weekStartsOnSunday bool
}

// NewGrumbleStatsUseCase creates a new instance.
func NewGrumbleStatsUseCase(repo grumble.Repository, defaultTZ string, weekStartsOnSunday bool) *GrumbleStatsUseCase {
	return &GrumbleStatsUseCase{repo: repo, defaultTZ: defaultTZ, weekStartsOnSunday: weekStartsOnSunday}
}

// StatsRequest holds parameters for stats retrieval.
type StatsRequest struct {
	Granularity grumble.Granularity
	From        *time.Time
	To          *time.Time
	TZ          string
}

// StatsByToxicRequest holds parameters for toxic-level stats retrieval.
type StatsByToxicRequest struct {
	StatsRequest
	ToxicLevel *int
}

// Get returns aggregated stats.
func (uc *GrumbleStatsUseCase) Get(ctx context.Context, req StatsRequest) ([]grumble.StatsRow, error) {
	fromUTC, toUTC, err := uc.resolveRange(req)
	if err != nil {
		return nil, err
	}

	return uc.repo.Stats(ctx, req.Granularity, fromUTC, toUTC)
}

// GetByToxic returns aggregated stats with toxic-level breakdown.
func (uc *GrumbleStatsUseCase) GetByToxic(ctx context.Context, req StatsByToxicRequest) ([]grumble.StatsRow, error) {
	if req.ToxicLevel != nil {
		if *req.ToxicLevel < 1 || *req.ToxicLevel > 5 {
			return nil, errors.New("toxic_level must be between 1 and 5")
		}
	}

	fromUTC, toUTC, err := uc.resolveRange(req.StatsRequest)
	if err != nil {
		return nil, err
	}

	return uc.repo.StatsByToxic(ctx, req.Granularity, fromUTC, toUTC, req.ToxicLevel)
}

func (uc *GrumbleStatsUseCase) resolveRange(req StatsRequest) (time.Time, time.Time, error) {
	if err := validateGranularity(req.Granularity); err != nil {
		return time.Time{}, time.Time{}, err
	}

	tz := req.TZ
	if tz == "" {
		tz = uc.defaultTZ
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if req.From != nil && req.To != nil {
		if !req.To.After(*req.From) {
			return time.Time{}, time.Time{}, errors.New("to must be after from")
		}
		return req.From.In(time.UTC), req.To.In(time.UTC), nil
	}

	start, end := uc.defaultRange(req.Granularity, loc)
	return start.In(time.UTC), end.In(time.UTC), nil
}

func validateGranularity(g grumble.Granularity) error {
	switch g {
	case grumble.GranularityDay, grumble.GranularityWeek, grumble.GranularityMonth:
		return nil
	default:
		return errors.New("invalid granularity")
	}
}

func (uc *GrumbleStatsUseCase) defaultRange(g grumble.Granularity, loc *time.Location) (time.Time, time.Time) {
	now := time.Now().In(loc)

	switch g {
	case grumble.GranularityDay:
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 0, 1)
	case grumble.GranularityWeek:
		start := startOfWeek(now, loc, uc.weekStartsOnSunday)
		return start, start.AddDate(0, 0, 7)
	case grumble.GranularityMonth:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0)
	default:
		return now, now
	}
}

func startOfWeek(t time.Time, loc *time.Location, sunday bool) time.Time {
	weekday := int(t.In(loc).Weekday()) // Sunday=0
	var offset int
	if sunday {
		offset = weekday
	} else {
		// Monday start: shift Sunday (0) to 6
		offset = (weekday + 6) % 7
	}
	start := t.AddDate(0, 0, -offset)
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
}
