package controller

import (
	"net/http"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
func (ctrl *TimelineController) GetGrumbles(
	c *gin.Context,
	userID *openapi_types.UUID,
	toxicLevelMin *int,
	toxicLevelMax *int,
	unpurifiedOnly *bool,
	limit *int,
	offset *int,
) {
	// Convert API params to domain types
	var toxicLevelMinDomain *shared.ToxicLevel
	var toxicLevelMaxDomain *shared.ToxicLevel
	var userIDFilter *shared.UserID

	if toxicLevelMin != nil {
		tl := shared.ToxicLevel(*toxicLevelMin)
		if err := tl.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": err.Error()})
			return
		}
		toxicLevelMinDomain = &tl
	}

	if toxicLevelMax != nil {
		tl := shared.ToxicLevel(*toxicLevelMax)
		if err := tl.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": err.Error()})
			return
		}
		toxicLevelMaxDomain = &tl
	}

	if toxicLevelMinDomain != nil && toxicLevelMaxDomain != nil && int(*toxicLevelMinDomain) > int(*toxicLevelMaxDomain) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY_PARAM", "message": "toxic_level_min cannot be greater than toxic_level_max"})
		return
	}

	if userID != nil {
		uid := shared.UserID(userID.String())
		userIDFilter = &uid
	}

	// Calculate page and page size
	page := 1
	pageSize := defaultTimelinePageSize
	offsetValue := 0

	if limit != nil && *limit > 0 {
		pageSize = *limit
	}
	if offset != nil && *offset >= 0 {
		offsetValue = *offset
	}
	if offsetValue > 0 {
		page = (offsetValue / pageSize) + 1
	}

	// Build usecase request
	req := usecase.TimelineRequest{
		ToxicLevelMin:  toxicLevelMinDomain,
		ToxicLevelMax:  toxicLevelMaxDomain,
		UnpurifiedOnly: unpurifiedOnly,
		UserID:         userIDFilter,
		Page:           page,
		PageSize:       pageSize,
		Offset:         offsetValue,
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
