package controller

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
)

// EventGrumblesController handles event grumbles retrieval
type EventGrumblesController struct {
	eventGrumblesGetUC *usecase.EventGrumblesGetUseCase
	grumblePresenter   *GrumblePresenter
	logger             logging.Logger
}

// NewEventGrumblesController creates a new EventGrumblesController
func NewEventGrumblesController(
	eventGrumblesGetUC *usecase.EventGrumblesGetUseCase,
	grumblePresenter *GrumblePresenter,
	logger logging.Logger,
) *EventGrumblesController {
	return &EventGrumblesController{
		eventGrumblesGetUC: eventGrumblesGetUC,
		grumblePresenter:   grumblePresenter,
		logger:             logger,
	}
}

// EventGrumblesQuery represents query parameters
type EventGrumblesQuery struct {
	ToxicLevelMin *int
	ToxicLevelMax *int
	Limit         int
	Offset        int
}

// EventGrumblesResponse represents the response
type EventGrumblesResponse struct {
	Grumbles      []*GrumbleResponse `json:"grumbles"`
	Total         int                `json:"total"`
	IsEventActive bool               `json:"is_event_active"`
	EventDate     string             `json:"event_date"` // YYYY-MM-DD形式
}

// GetEventGrumbles retrieves event grumbles from archive
func (ctrl *EventGrumblesController) GetEventGrumbles(
	ctx context.Context,
	query EventGrumblesQuery,
) (*EventGrumblesResponse, error) {
	var toxicLevelMin *shared.ToxicLevel
	var toxicLevelMax *shared.ToxicLevel

	if query.ToxicLevelMin != nil {
		tl := shared.ToxicLevel(*query.ToxicLevelMin)
		toxicLevelMin = &tl
	}

	if query.ToxicLevelMax != nil {
		tl := shared.ToxicLevel(*query.ToxicLevelMax)
		toxicLevelMax = &tl
	}

	req := usecase.EventGrumblesRequest{
		ToxicLevelMin:   toxicLevelMin,
		ToxicLevelMax:   toxicLevelMax,
		ExcludePurified: false, // イベントでは全て表示
		Limit:           query.Limit,
		Offset:          query.Offset,
	}

	result, err := ctrl.eventGrumblesGetUC.Get(ctx, req)
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to get event grumbles", "error", err)
		return nil, err
	}

	apiGrumbles, err := ctrl.grumblePresenter.ToAPIGrumbles(result.Grumbles)
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to convert grumbles", "error", err)
		return nil, err
	}

	eventDateStr := ""
	if result.IsEventActive {
		eventDateStr = result.EventDate.Format("2006-01-02")
	}

	return &EventGrumblesResponse{
		Grumbles:      apiGrumbles,
		Total:         result.Total,
		IsEventActive: result.IsEventActive,
		EventDate:     eventDateStr,
	}, nil
}
