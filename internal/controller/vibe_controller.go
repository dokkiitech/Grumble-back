package controller

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/google/uuid"
)

// VibeController handles vibe-related application logic.
type VibeController struct {
	vibeAddUC *usecase.VibeAddUseCase
	logger    logging.Logger
}

// NewVibeController creates a new VibeController.
func NewVibeController(
	vibeAddUC *usecase.VibeAddUseCase,
	logger logging.Logger,
) *VibeController {
	return &VibeController{
		vibeAddUC: vibeAddUC,
		logger:    logger,
	}
}

// AddVibeInput represents the application-level request.
type AddVibeInput struct {
	GrumbleID shared.GrumbleID
	UserID    shared.UserID
	VibeType  shared.VibeType
}

// AddVibeResponse represents the response returned to the HTTP layer.
type AddVibeResponse struct {
	VibeID       int
	GrumbleID    uuid.UUID
	UserID       uuid.UUID
	VibeType     shared.VibeType
	VotedAt      time.Time
	VibeCount    int
	VirtuePoints int
	IsPurified   bool
}

// AddVibe executes the use case and returns the response model.
func (ctrl *VibeController) AddVibe(ctx context.Context, input AddVibeInput) (*AddVibeResponse, error) {
	vibeType := input.VibeType
	if vibeType == "" {
		vibeType = shared.VibeTypeWakaru
	}

	ucReq := usecase.AddVibeRequest{
		GrumbleID: input.GrumbleID,
		UserID:    input.UserID,
		VibeType:  vibeType,
	}

	resp, err := ctrl.vibeAddUC.Add(ctx, ucReq)
	if err != nil {
		return nil, err
	}

	grumbleUUID, err := uuid.Parse(string(resp.Vibe.GrumbleID))
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to parse grumble UUID", "error", err)
		return nil, err
	}

	userUUID, err := uuid.Parse(string(resp.Vibe.UserID))
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to parse user UUID", "error", err)
		return nil, err
	}

	return &AddVibeResponse{
		VibeID:       int(resp.Vibe.VibeID),
		GrumbleID:    grumbleUUID,
		UserID:       userUUID,
		VibeType:     resp.Vibe.Type,
		VotedAt:      resp.Vibe.VotedAt,
		VibeCount:    resp.VibeCount,
		VirtuePoints: resp.VirtuePoints,
		IsPurified:   resp.IsPurified,
	}, nil
}
