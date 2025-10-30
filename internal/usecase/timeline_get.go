package usecase

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
)

// TimelineGetUseCase handles retrieving the timeline of grumbles
type TimelineGetUseCase struct {
	grumbleRepo grumble.Repository
}

// NewTimelineGetUseCase creates a new TimelineGetUseCase
func NewTimelineGetUseCase(grumbleRepo grumble.Repository) *TimelineGetUseCase {
	return &TimelineGetUseCase{
		grumbleRepo: grumbleRepo,
	}
}

// TimelineRequest represents the input for getting timeline
type TimelineRequest struct {
	ToxicLevelMin  *shared.ToxicLevel // Optional minimum toxic level (inclusive)
	ToxicLevelMax  *shared.ToxicLevel // Optional maximum toxic level (inclusive)
	UnpurifiedOnly *bool              // Optional flag to restrict to unpurified grumbles
	Page           int                // Page number (1-indexed)
	PageSize       int                // Number of items per page
	Offset         int                // Number of items to skip before starting results
}

// TimelineResponse represents the timeline result
type TimelineResponse struct {
	Grumbles   []*grumble.Grumble
	TotalCount int
	Page       int
	PageSize   int
}

// Get retrieves the timeline with filters
func (uc *TimelineGetUseCase) Get(ctx context.Context, req TimelineRequest) (*TimelineResponse, error) {
	// Set defaults
	if req.PageSize <= 0 {
		req.PageSize = 20 // Default page size
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}
	if offset > 0 {
		req.Page = (offset / req.PageSize) + 1
	} else {
		if req.Page <= 0 {
			req.Page = 1
		}
		offset = (req.Page - 1) * req.PageSize
	}

	// Build filter
	excludePurified := true
	if req.UnpurifiedOnly != nil {
		excludePurified = *req.UnpurifiedOnly
	}

	filter := grumble.TimelineFilter{
		ToxicLevelMin:   req.ToxicLevelMin,
		ToxicLevelMax:   req.ToxicLevelMax,
		ExcludePurified: excludePurified,
		ExcludeExpired:  true, // Always exclude expired grumbles
		Limit:           req.PageSize,
		Offset:          offset,
	}

	// Get grumbles
	grumbles, err := uc.grumbleRepo.FindTimeline(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	totalCount, err := uc.grumbleRepo.CountTimeline(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &TimelineResponse{
		Grumbles:   grumbles,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}
