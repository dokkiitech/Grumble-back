package controller

import (
	"net/http"
	"strconv"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const defaultTimelinePageSize = 20

// TimelineController handles timeline-related HTTP requests
type TimelineController struct {
	timelineGetUC *usecase.TimelineGetUseCase
	presenter     *TimelinePresenter
	logger        logging.Logger
}

// NewTimelineController creates a new TimelineController
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

// GetGrumbles handles GET /grumbles (�����֗)
func (ctrl *TimelineController) GetGrumbles(c *gin.Context) {
	// Parse query parameters
	var toxicLevelMin *shared.ToxicLevel
	var toxicLevelMax *shared.ToxicLevel
	var unpurifiedOnly *bool
	var userIDFilter *shared.UserID

	if v := c.Query("toxic_level_min"); v != "" {
		iv, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": "toxic_level_min must be an integer between 1 and 5"})
			return
		}
		tl := shared.ToxicLevel(iv)
		if err := tl.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": err.Error()})
			return
		}
		toxicLevelMin = &tl
	}

	if v := c.Query("toxic_level_max"); v != "" {
		iv, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": "toxic_level_max must be an integer between 1 and 5"})
			return
		}
		tl := shared.ToxicLevel(iv)
		if err := tl.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": err.Error()})
			return
		}
		toxicLevelMax = &tl
	}

	if toxicLevelMin != nil && toxicLevelMax != nil && int(*toxicLevelMin) > int(*toxicLevelMax) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": "toxic_level_min cannot be greater than toxic_level_max"})
		return
	}

	if v := c.Query("unpurified_only"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			unpurifiedOnly = &b
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": "unpurified_only must be a boolean"})
			return
		}
	}

	if v := c.Query("user_id"); v != "" {
		if _, err := uuid.Parse(v); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": "user_id must be a valid UUID"})
			return
		}
		uid := shared.UserID(v)
		userIDFilter = &uid
	}

	// Calculate page and page size
	page := 1
	pageSize := defaultTimelinePageSize
	offset := 0

	if v := c.Query("limit"); v != "" {
		if iv, err := strconv.Atoi(v); err == nil && iv > 0 {
			pageSize = iv
		}
	}
	if v := c.Query("offset"); v != "" {
		if iv, err := strconv.Atoi(v); err == nil && iv >= 0 {
			offset = iv
		}
	}
	if offset > 0 {
		page = (offset / pageSize) + 1
	}

	// Build usecase request
	req := usecase.TimelineRequest{
		ToxicLevelMin:  toxicLevelMin,
		ToxicLevelMax:  toxicLevelMax,
		UnpurifiedOnly: unpurifiedOnly,
		UserID:         userIDFilter,
		Page:           page,
		PageSize:       pageSize,
		Offset:         offset,
	}

	// Execute use case
	resp, err := ctrl.timelineGetUC.Get(c.Request.Context(), req)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to get timeline", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to retrieve timeline"})
		return
	}

	// Convert to API response
	apiResp, err := ctrl.presenter.ToAPITimelineResponse(resp)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to convert timeline to API response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "message": "Failed to format response"})
		return
	}
	c.JSON(http.StatusOK, apiResp)
}
