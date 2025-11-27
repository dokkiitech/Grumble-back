package usecase

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	sharedservice "github.com/dokkiitech/grumble-back/internal/domain/shared/service"
	"github.com/google/uuid"
)

// GrumblePostUseCase handles posting new grumbles
type GrumblePostUseCase struct {
	grumbleRepo  grumble.Repository
	eventTimeSvc *sharedservice.EventTimeService
}

// NewGrumblePostUseCase creates a new GrumblePostUseCase
func NewGrumblePostUseCase(grumbleRepo grumble.Repository, eventTimeSvc *sharedservice.EventTimeService) *GrumblePostUseCase {
	return &GrumblePostUseCase{
		grumbleRepo:  grumbleRepo,
		eventTimeSvc: eventTimeSvc,
	}
}

// PostGrumbleRequest represents the input for posting a grumble
type PostGrumbleRequest struct {
	UserID         shared.UserID
	Content        string
	ToxicLevel     shared.ToxicLevel
	IsEventGrumble bool
}

// Post creates and persists a new grumble
func (uc *GrumblePostUseCase) Post(ctx context.Context, req PostGrumbleRequest) (*grumble.Grumble, error) {
	// Create grumble entity
	now := time.Now()
	g := &grumble.Grumble{
		GrumbleID:      shared.GrumbleID(uuid.New().String()),
		UserID:         req.UserID,
		Content:        req.Content,
		ToxicLevel:     req.ToxicLevel,
		VibeCount:      0,
		IsPurified:     false,
		PostedAt:       now,
		ExpiresAt:      uc.eventTimeSvc.CalculateNextMidnight(now), // 翌日の00:00
		IsEventGrumble: req.IsEventGrumble,
	}

	// Validate business rules
	if err := g.Validate(); err != nil {
		return nil, err
	}

	// Persist to repository
	if err := uc.grumbleRepo.Create(ctx, g); err != nil {
		return nil, err
	}

	return g, nil
}
