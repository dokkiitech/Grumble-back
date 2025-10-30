package controller

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VibeController handles vibe-related HTTP requests.
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

// PostVibe handles POST /grumbles/{grumble_id}/vibes.
func (ctrl *VibeController) PostVibe(c *gin.Context, grumbleID uuid.UUID) {
	type requestBody struct {
		VibeType *string `json:"vibe_type,omitempty"`
	}

	var body requestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		if !errors.Is(err, io.EOF) {
			ctrl.logger.WarnContext(c.Request.Context(), "Invalid AddVibe request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "message": "Invalid request body"})
			return
		}
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": "User not authenticated"})
		return
	}
	userID, ok := userIDValue.(shared.UserID)
	if !ok {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid user_id type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to authenticate user"})
		return
	}

	vibeType := shared.VibeTypeWakaru
	if body.VibeType != nil {
		vibeType = shared.VibeType(strings.ToUpper(strings.TrimSpace(*body.VibeType)))
	}

	ucReq := usecase.AddVibeRequest{
		GrumbleID: shared.GrumbleID(grumbleID.String()),
		UserID:    userID,
		VibeType:  vibeType,
	}

	resp, err := ctrl.vibeAddUC.Add(c.Request.Context(), ucReq)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	userUUIDParsed, err := uuid.Parse(string(resp.Vibe.UserID))
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to parse user UUID for response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to format response"})
		return
	}

	response := gin.H{
		"vibe_id":       int(resp.Vibe.VibeID),
		"grumble_id":    grumbleID,
		"user_id":       userUUIDParsed,
		"vibe_type":     string(resp.Vibe.Type),
		"voted_at":      resp.Vibe.VotedAt,
		"vibe_count":    resp.VibeCount,
		"virtue_points": resp.VirtuePoints,
		"is_purified":   resp.IsPurified,
	}

	c.JSON(http.StatusCreated, response)
}

func (ctrl *VibeController) handleError(c *gin.Context, err error) {
	var validationErr *shared.ValidationError
	var notFoundErr *shared.NotFoundError
	var duplicateErr *shared.DuplicateVibeError
	var unauthorizedErr *shared.UnauthorizedError
	var internalErr *shared.InternalError

	switch {
	case errors.As(err, &validationErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": validationErr.Error()})
	case errors.As(err, &notFoundErr):
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": notFoundErr.Error()})
	case errors.As(err, &duplicateErr):
		c.JSON(http.StatusConflict, gin.H{"error": "DUPLICATE_VIBE", "message": duplicateErr.Error()})
	case errors.As(err, &unauthorizedErr):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": unauthorizedErr.Error()})
	case errors.As(err, &internalErr):
		ctrl.logger.ErrorContext(c.Request.Context(), "Internal error during AddVibe", "error", internalErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "An internal error occurred"})
	default:
		ctrl.logger.ErrorContext(c.Request.Context(), "Unknown error during AddVibe", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UNKNOWN_ERROR", "message": "An unknown error occurred"})
	}
}
