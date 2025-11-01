package controller

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
)

const defaultTimelinePageSize = 20

// TimelineController handles timeline-related application logic.
type TimelineController struct {
	timelineGetUC *usecase.TimelineGetUseCase
	presenter     *TimelinePresenter
	logger        logging.Logger
}

// NewTimelineController creates a new TimelineController.
func NewTimelineController(
	timelineGetUC *usecase.TimelineGetUseCase,
	presenter *TimelinePresenter,
	logger logging.Logger,
) *TimelineController {
	return &TimelineController{
		timelineGetUC: timelineGetUC,
		presenter:     presenter,
		logger:        logger,
	}
}

// TimelineQuery represents filters supplied by the HTTP layer.
type TimelineQuery struct {
	UserID        *shared.UserID
	ViewerUserID  *shared.UserID
	ToxicLevelMin *shared.ToxicLevel
	ToxicLevelMax *shared.ToxicLevel
	IsPurified    *bool
	Limit         int
	Offset        int
}

// GetGrumbles retrieves the timeline and returns the API-facing response model.
func (ctrl *TimelineController) GetGrumbles(ctx context.Context, query TimelineQuery) (*TimelineResponse, error) {
	pageSize := query.Limit
	if pageSize <= 0 {
		pageSize = defaultTimelinePageSize
	}

	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	page := 1
	if offset > 0 {
		page = (offset / pageSize) + 1
	}

	req := usecase.TimelineRequest{
		ToxicLevelMin: query.ToxicLevelMin,
		ToxicLevelMax: query.ToxicLevelMax,
		IsPurified:    query.IsPurified,
		UserID:        query.UserID,
		ViewerUserID:  query.ViewerUserID,
		Page:          page,
		PageSize:      pageSize,
		Offset:        offset,
	}

	resp, err := ctrl.timelineGetUC.Get(ctx, req)
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to get timeline", "error", err)
		return nil, err
	}

	apiResp, err := ctrl.presenter.ToAPITimelineResponse(resp)
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to convert timeline to API response", "error", err)
		return nil, err
	}

	return apiResp, nil
}
