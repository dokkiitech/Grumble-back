package controller

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
)

// GrumbleController handles grumble-related application logic.
type GrumbleController struct {
	postGrumbleUC *usecase.GrumblePostUseCase
	presenter     *GrumblePresenter
	logger        logging.Logger
}

// NewGrumbleController creates a new GrumbleController.
func NewGrumbleController(
	postGrumbleUC *usecase.GrumblePostUseCase,
	presenter *GrumblePresenter,
	logger logging.Logger,
) *GrumbleController {
	return &GrumbleController{
		postGrumbleUC: postGrumbleUC,
		presenter:     presenter,
		logger:        logger,
	}
}

// CreateGrumbleInput is the application-level request for creating a grumble.
type CreateGrumbleInput struct {
	UserID         shared.UserID
	Content        string
	ToxicLevel     shared.ToxicLevel
	IsEventGrumble bool
}

// CreateGrumble executes the use case and returns the API-facing response model.
func (ctrl *GrumbleController) CreateGrumble(ctx context.Context, input CreateGrumbleInput) (*GrumbleResponse, error) {
	ucReq := usecase.PostGrumbleRequest{
		UserID:         input.UserID,
		Content:        input.Content,
		ToxicLevel:     input.ToxicLevel,
		IsEventGrumble: input.IsEventGrumble,
	}

	grumble, err := ctrl.postGrumbleUC.Post(ctx, ucReq)
	if err != nil {
		return nil, err
	}

	response, err := ctrl.presenter.ToAPIGrumble(grumble)
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to convert grumble to API response", "error", err)
		return nil, err
	}

	return response, nil
}
