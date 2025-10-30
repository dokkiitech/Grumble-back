package controller

import (
	"errors"
	"net/http"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GrumbleController handles grumble-related HTTP requests
type GrumbleController struct {
	postGrumbleUC *usecase.GrumblePostUseCase
	presenter     *GrumblePresenter
	logger        logging.Logger
}

// NewGrumbleController creates a new GrumbleController
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

// CreateGrumble handles POST /grumbles (�?\)
func (ctrl *GrumbleController) CreateGrumble(c *gin.Context) {
	type createGrumbleRequest struct {
		Content        string `json:"content"`
		IsEventGrumble *bool  `json:"is_event_grumble,omitempty"`
		ToxicLevel     int    `json:"toxic_level"`
	}
	type errorResponse struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	// Parse request body
	var req createGrumbleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, errorResponse{Error: "INVALID_REQUEST", Message: "Invalid request body"})
		return
	}

	// Get user ID from context (set by auth middleware)
	// For now, use a placeholder - will be implemented in US3
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": "User not authenticated"})
		return
	}
	userID := userIDValue.(shared.UserID)

	// Convert to usecase request
	isEventGrumble := false
	if req.IsEventGrumble != nil {
		isEventGrumble = *req.IsEventGrumble
	}

	ucReq := usecase.PostGrumbleRequest{
		UserID:         userID,
		Content:        req.Content,
		ToxicLevel:     shared.ToxicLevel(req.ToxicLevel),
		IsEventGrumble: isEventGrumble,
	}

	// Execute use case
	grumble, err := ctrl.postGrumbleUC.Post(c.Request.Context(), ucReq)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	// Convert to API response
	response, err := ctrl.presenter.ToAPIGrumble(grumble)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to convert grumble to API response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to format response"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// handleError maps domain errors to HTTP status codes
func (ctrl *GrumbleController) handleError(c *gin.Context, err error) {
	var validationErr *shared.ValidationError
	var notFoundErr *shared.NotFoundError
	var unauthorizedErr *shared.UnauthorizedError
	var internalErr *shared.InternalError

	switch {
	case errors.As(err, &validationErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": validationErr.Error()})
	case errors.As(err, &notFoundErr):
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": notFoundErr.Error()})
	case errors.As(err, &unauthorizedErr):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": unauthorizedErr.Error()})
	case errors.As(err, &internalErr):
		ctrl.logger.ErrorContext(c.Request.Context(), "Internal error", "error", internalErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "An internal error occurred"})
	default:
		ctrl.logger.ErrorContext(c.Request.Context(), "Unknown error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UNKNOWN_ERROR", "message": "An unknown error occurred"})
	}
}
