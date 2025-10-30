package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authAnonymousUC *usecase.AuthAnonymousUseCase
	userQueryUC     *usecase.UserQueryUseCase
	logger          logging.Logger
}

// NewAuthController creates a new AuthController
func NewAuthController(
	authAnonymousUC *usecase.AuthAnonymousUseCase,
	userQueryUC *usecase.UserQueryUseCase,
	logger logging.Logger,
) *AuthController {
	return &AuthController{
		authAnonymousUC: authAnonymousUC,
		userQueryUC:     userQueryUC,
		logger:          logger,
	}
}

// GetMyProfile handles GET /users/me (�n�����1֗)
func (ctrl *AuthController) GetMyProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": "User not authenticated"})
		return
	}
	userID := userIDValue.(shared.UserID)

	// Get user from repository
	user, err := ctrl.userQueryUC.GetMyProfile(c.Request.Context(), userID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	// Build response
	userUUID, err := uuid.Parse(string(user.UserID))
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to parse user UUID", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to format response"})
		return
	}

	resp := gin.H{
		"user_id":       userUUID,
		"virtue_points": user.VirtuePoints,
		"created_at":    user.CreatedAt,
	}
	if user.ProfileTitle != nil {
		resp["profile_title"] = *user.ProfileTitle
	}

	c.JSON(http.StatusOK, resp)
}

// GetBodhisattvaRankings handles GET /users/rankings (����֗)
// Note: This endpoint may not be in the OpenAPI spec yet
func (ctrl *AuthController) GetBodhisattvaRankings(c *gin.Context) {
	// Get limit from query param, default to 10
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get top users
	users, err := ctrl.userQueryUC.GetBodhisattvaRankings(c.Request.Context(), limit)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	// Convert to API response
	rankings := make([]gin.H, len(users))
	for i, u := range users {
		userUUID, err := uuid.Parse(string(u.UserID))
		if err != nil {
			ctrl.logger.ErrorContext(c.Request.Context(), "Failed to parse user UUID", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to format response"})
			return
		}
		m := gin.H{"user_id": userUUID, "virtue_points": u.VirtuePoints, "created_at": u.CreatedAt}
		if u.ProfileTitle != nil {
			m["profile_title"] = *u.ProfileTitle
		}
		rankings[i] = m
	}

	c.JSON(http.StatusOK, gin.H{"rankings": rankings, "total": len(rankings)})
}

// handleError maps domain errors to HTTP status codes
func (ctrl *AuthController) handleError(c *gin.Context, err error) {
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
