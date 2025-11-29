package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dokkiitech/grumble-back/internal/controller"
	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/nullable"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// StrictControllerServer bridges HTTP layer and application controllers using the strict server interface.
type StrictControllerServer struct {
	grumbleController       *controller.GrumbleController
	timelineController      *controller.TimelineController
	authController          *controller.AuthController
	vibeController          *controller.VibeController
	eventGrumblesController *controller.EventGrumblesController
	statsController         *controller.GrumbleStatsController
	logger                  logging.Logger
}

// NewStrictControllerServer constructs a StrictServerInterface backed by existing controllers.
func NewStrictControllerServer(
	grumbleCtrl *controller.GrumbleController,
	timelineCtrl *controller.TimelineController,
	authCtrl *controller.AuthController,
	vibeCtrl *controller.VibeController,
	eventGrumblesCtrl *controller.EventGrumblesController,
	statsCtrl *controller.GrumbleStatsController,
	logger logging.Logger,
) *StrictControllerServer {
	return &StrictControllerServer{
		grumbleController:       grumbleCtrl,
		timelineController:      timelineCtrl,
		authController:          authCtrl,
		vibeController:          vibeCtrl,
		eventGrumblesController: eventGrumblesCtrl,
		statsController:         statsCtrl,
		logger:                  logger,
	}
}

// GetEvents is currently not implemented.
func (s *StrictControllerServer) GetEvents(ctx context.Context, _ GetEventsRequestObject) (GetEventsResponseObject, error) {
	return nil, errors.New("GetEvents not implemented")
}

// GetEvent is currently not implemented.
func (s *StrictControllerServer) GetEvent(ctx context.Context, _ GetEventRequestObject) (GetEventResponseObject, error) {
	return nil, errors.New("GetEvent not implemented")
}

// GetGrumbleStats handles aggregated stats retrieval.
func (s *StrictControllerServer) GetGrumbleStats(ctx context.Context, request GetGrumbleStatsRequestObject) (GetGrumbleStatsResponseObject, error) {
	params := request.Params

	tz := ""
	if params.Tz != nil {
		tz = *params.Tz
	}

	input := controller.StatsInput{
		Granularity: grumble.Granularity(params.Granularity),
		TZ:          tz,
		From:        params.From,
		To:          params.To,
	}

	result, err := s.statsController.GetStats(ctx, input)
	if err != nil {
		s.logger.ErrorContext(ctx, "GetGrumbleStats failed", "error", err)
		if isValidationError(err) {
			return GetGrumbleStats400JSONResponse(errorResponse("INVALID_QUERY_PARAM", err.Error())), nil
		}
		return GetGrumbleStats401JSONResponse(errorResponse("INTERNAL_ERROR", "Failed to retrieve stats")), nil
	}

	response := make([]GrumbleStatsBucket, len(result))
	for i, r := range result {
		response[i] = GrumbleStatsBucket{
			Bucket:          r.Bucket,
			PurifiedCount:   r.PurifiedCount,
			UnpurifiedCount: r.UnpurifiedCount,
			TotalVibes:      r.TotalVibes,
		}
	}

	return GetGrumbleStats200JSONResponse(response), nil
}

// GetGrumbleStatsToxic handles aggregated stats with toxic-level breakdown.
func (s *StrictControllerServer) GetGrumbleStatsToxic(ctx context.Context, request GetGrumbleStatsToxicRequestObject) (GetGrumbleStatsToxicResponseObject, error) {
	params := request.Params

	tz := ""
	if params.Tz != nil {
		tz = *params.Tz
	}

	input := controller.StatsByToxicInput{
		StatsInput: controller.StatsInput{
			Granularity: grumble.Granularity(params.Granularity),
			TZ:          tz,
			From:        params.From,
			To:          params.To,
		},
		ToxicLevel: params.ToxicLevel,
	}

	result, err := s.statsController.GetStatsByToxic(ctx, input)
	if err != nil {
		s.logger.ErrorContext(ctx, "GetGrumbleStatsToxic failed", "error", err)
		if isValidationError(err) {
			return GetGrumbleStatsToxic400JSONResponse(errorResponse("INVALID_QUERY_PARAM", err.Error())), nil
		}
		return GetGrumbleStatsToxic401JSONResponse(errorResponse("INTERNAL_ERROR", "Failed to retrieve stats")), nil
	}

	response := make([]GrumbleStatsToxicBucket, len(result))
	for i, r := range result {
		level := 0
		if r.ToxicLevel != nil {
			level = *r.ToxicLevel
		}
		response[i] = GrumbleStatsToxicBucket{
			Bucket:          r.Bucket,
			ToxicLevel:      level,
			PurifiedCount:   r.PurifiedCount,
			UnpurifiedCount: r.UnpurifiedCount,
			TotalVibes:      r.TotalVibes,
		}
	}

	return GetGrumbleStatsToxic200JSONResponse(response), nil
}

// GetEventGrumbles handles retrieving event grumbles from archive
func (s *StrictControllerServer) GetEventGrumbles(ctx context.Context, request GetEventGrumblesRequestObject) (GetEventGrumblesResponseObject, error) {
	params := request.Params

	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	query := controller.EventGrumblesQuery{
		ToxicLevelMin: params.ToxicLevelMin,
		ToxicLevelMax: params.ToxicLevelMax,
		Limit:         limit,
		Offset:        offset,
	}

	response, err := s.eventGrumblesController.GetEventGrumbles(ctx, query)
	if err != nil {
		s.logger.ErrorContext(ctx, "GetEventGrumbles failed", "error", err)
		return GetEventGrumbles401JSONResponse(errorResponse("INTERNAL_ERROR", "Failed to retrieve event grumbles")), nil
	}

	// Convert controller response to OpenAPI response
	apiGrumbles := make([]Grumble, len(response.Grumbles))
	for i, g := range response.Grumbles {
		apiGrumbles[i] = toAPIGrumble(g)
	}

	eventDate := openapi_types.Date{}
	if response.EventDate != "" {
		parsedTime, err := time.Parse("2006-01-02", response.EventDate)
		if err == nil {
			eventDate = openapi_types.Date{Time: parsedTime}
		}
	}

	return GetEventGrumbles200JSONResponse{
		Grumbles:      apiGrumbles,
		Total:         response.Total,
		IsEventActive: response.IsEventActive,
		EventDate:     eventDate,
	}, nil
}

// GetGrumbles handles the timeline retrieval.
func (s *StrictControllerServer) GetGrumbles(ctx context.Context, request GetGrumblesRequestObject) (GetGrumblesResponseObject, error) {
	params := request.Params

	var toxicLevelMin *shared.ToxicLevel
	if params.ToxicLevelMin != nil {
		tl := shared.ToxicLevel(*params.ToxicLevelMin)
		if err := tl.Validate(); err != nil {
			return GetGrumbles400JSONResponse(errorResponse("INVALID_QUERY_PARAM", err.Error())), nil
		}
		toxicLevelMin = &tl
	}

	var toxicLevelMax *shared.ToxicLevel
	if params.ToxicLevelMax != nil {
		tl := shared.ToxicLevel(*params.ToxicLevelMax)
		if err := tl.Validate(); err != nil {
			return GetGrumbles400JSONResponse(errorResponse("INVALID_QUERY_PARAM", err.Error())), nil
		}
		toxicLevelMax = &tl
	}

	if toxicLevelMin != nil && toxicLevelMax != nil && int(*toxicLevelMin) > int(*toxicLevelMax) {
		return GetGrumbles400JSONResponse(errorResponse("INVALID_QUERY_PARAM", "toxic_level_min cannot be greater than toxic_level_max")), nil
	}

	var userID *shared.UserID
	if params.UserID != nil {
		uid := shared.UserID(params.UserID.String())
		userID = &uid
	}

	limit := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	var viewerUserID *shared.UserID
	if id, ok := s.userIDFromContext(ctx); ok {
		viewer := id
		viewerUserID = &viewer
	}

	query := controller.TimelineQuery{
		UserID:        userID,
		ViewerUserID:  viewerUserID,
		ToxicLevelMin: toxicLevelMin,
		ToxicLevelMax: toxicLevelMax,
		IsPurified:    params.IsPurified,
		Limit:         limit,
		Offset:        offset,
	}

	result, err := s.timelineController.GetGrumbles(ctx, query)
	if err != nil {
		if resp, ok := s.timelineErrorResponse(ctx, err); ok {
			return resp, nil
		}
		return nil, err
	}

	grumbles := make([]Grumble, len(result.Grumbles))
	for i, g := range result.Grumbles {
		grumbles[i] = toAPIGrumble(g)
	}
	total := result.Total

	return GetGrumbles200JSONResponse{Grumbles: &grumbles, Total: &total}, nil
}

// CreateGrumble handles POST /grumbles.
func (s *StrictControllerServer) CreateGrumble(ctx context.Context, request CreateGrumbleRequestObject) (CreateGrumbleResponseObject, error) {
	if request.Body == nil {
		return CreateGrumble400JSONResponse(errorResponse("INVALID_REQUEST", "request body is required")), nil
	}

	userID, ok := s.userIDFromContext(ctx)
	if !ok {
		return CreateGrumble401JSONResponse(errorResponse("UNAUTHORIZED", "User not authenticated")), nil
	}

	toxicLevel := shared.ToxicLevel(request.Body.ToxicLevel)
	if err := toxicLevel.Validate(); err != nil {
		return CreateGrumble400JSONResponse(errorResponse("VALIDATION_ERROR", err.Error())), nil
	}

	input := controller.CreateGrumbleInput{
		UserID:            userID,
		Content:           request.Body.Content,
		ToxicLevel:        toxicLevel,
		PurifiedThreshold: request.Body.PurifiedThreshold,
		IsEventGrumble:    request.Body.IsEventGrumble != nil && *request.Body.IsEventGrumble,
	}

	grumble, err := s.grumbleController.CreateGrumble(ctx, input)
	if err != nil {
		if resp, ok := s.createGrumbleErrorResponse(ctx, err); ok {
			return resp, nil
		}
		return nil, err
	}

	apiGrumble := toAPIGrumble(grumble)
	return CreateGrumble201JSONResponse(apiGrumble), nil
}

// AddVibe handles POST /grumbles/{grumble_id}/vibes.
func (s *StrictControllerServer) AddVibe(ctx context.Context, request AddVibeRequestObject) (AddVibeResponseObject, error) {
	userID, ok := s.userIDFromContext(ctx)
	if !ok {
		return AddVibe401JSONResponse(errorResponse("UNAUTHORIZED", "User not authenticated")), nil
	}

	vibeType := shared.VibeTypeWakaru
	if request.Body != nil && request.Body.VibeType != nil {
		vibeType = shared.VibeType(strings.ToUpper(strings.TrimSpace(string(*request.Body.VibeType))))
	}

	if err := vibeType.Validate(); err != nil {
		return AddVibe400JSONResponse(errorResponse("VALIDATION_ERROR", err.Error())), nil
	}

	input := controller.AddVibeInput{
		GrumbleID: shared.GrumbleID(request.GrumbleID.String()),
		UserID:    userID,
		VibeType:  vibeType,
	}

	vibe, err := s.vibeController.AddVibe(ctx, input)
	if err != nil {
		if resp, ok := s.addVibeErrorResponse(ctx, err); ok {
			return resp, nil
		}
		return nil, err
	}

	apiVibe := toAPIVibe(vibe)
	return AddVibe201JSONResponse(apiVibe), nil
}

// GetMyProfile handles GET /users/me.
func (s *StrictControllerServer) GetMyProfile(ctx context.Context, _ GetMyProfileRequestObject) (GetMyProfileResponseObject, error) {
	userID, ok := s.userIDFromContext(ctx)
	if !ok {
		return GetMyProfile401JSONResponse(errorResponse("UNAUTHORIZED", "User not authenticated")), nil
	}

	profile, err := s.authController.GetMyProfile(ctx, userID)
	if err != nil {
		if resp, ok := s.profileErrorResponse(ctx, err); ok {
			return resp, nil
		}
		return nil, err
	}

	apiProfile := toAPIAnonymousUser(profile)
	return GetMyProfile200JSONResponse(apiProfile), nil
}

func (s *StrictControllerServer) timelineErrorResponse(ctx context.Context, err error) (GetGrumblesResponseObject, bool) {
	if classification, ok := s.classifyError(ctx, err); ok {
		switch classification.Status {
		case http.StatusBadRequest:
			return GetGrumbles400JSONResponse(classification.Payload), true
		case http.StatusUnauthorized:
			return GetGrumbles401JSONResponse(classification.Payload), true
		}
	}
	return nil, false
}

func (s *StrictControllerServer) createGrumbleErrorResponse(ctx context.Context, err error) (CreateGrumbleResponseObject, bool) {
	if classification, ok := s.classifyError(ctx, err); ok {
		switch classification.Status {
		case http.StatusBadRequest:
			return CreateGrumble400JSONResponse(classification.Payload), true
		case http.StatusUnauthorized:
			return CreateGrumble401JSONResponse(classification.Payload), true
		}
	}
	return nil, false
}

func (s *StrictControllerServer) addVibeErrorResponse(ctx context.Context, err error) (AddVibeResponseObject, bool) {
	if classification, ok := s.classifyError(ctx, err); ok {
		switch classification.Status {
		case http.StatusBadRequest:
			return AddVibe400JSONResponse(classification.Payload), true
		case http.StatusUnauthorized:
			return AddVibe401JSONResponse(classification.Payload), true
		case http.StatusNotFound:
			return AddVibe404JSONResponse(classification.Payload), true
		case http.StatusConflict:
			return AddVibe409JSONResponse(classification.Payload), true
		}
	}
	return nil, false
}

func (s *StrictControllerServer) profileErrorResponse(ctx context.Context, err error) (GetMyProfileResponseObject, bool) {
	if classification, ok := s.classifyError(ctx, err); ok {
		switch classification.Status {
		case http.StatusUnauthorized:
			return GetMyProfile401JSONResponse(classification.Payload), true
		}
	}
	return nil, false
}

type errorClassification struct {
	Status  int
	Payload ErrorResponse
}

func (s *StrictControllerServer) classifyError(ctx context.Context, err error) (errorClassification, bool) {
	var (
		validationErr           *shared.ValidationError
		notFoundErr             *shared.NotFoundError
		duplicateErr            *shared.DuplicateVibeError
		unauthorizedErr         *shared.UnauthorizedError
		inappropriateContentErr *shared.InappropriateContentError
		internalErr             *shared.InternalError
	)

	switch {
	case errors.As(err, &validationErr):
		return errorClassification{Status: http.StatusBadRequest, Payload: errorResponse("VALIDATION_ERROR", validationErr.Error())}, true
	case errors.As(err, &inappropriateContentErr):
		return errorClassification{Status: http.StatusBadRequest, Payload: errorResponse("INAPPROPRIATE_CONTENT", inappropriateContentErr.Error())}, true
	case errors.As(err, &notFoundErr):
		return errorClassification{Status: http.StatusNotFound, Payload: errorResponse("NOT_FOUND", notFoundErr.Error())}, true
	case errors.As(err, &duplicateErr):
		return errorClassification{Status: http.StatusConflict, Payload: errorResponse("DUPLICATE_VIBE", duplicateErr.Error())}, true
	case errors.As(err, &unauthorizedErr):
		return errorClassification{Status: http.StatusUnauthorized, Payload: errorResponse("UNAUTHORIZED", unauthorizedErr.Error())}, true
	case errors.As(err, &internalErr):
		s.logger.ErrorContext(ctx, "Internal error", "error", internalErr)
		return errorClassification{}, false
	default:
		s.logger.ErrorContext(ctx, "Unhandled error", "error", err)
		return errorClassification{}, false
	}
}

func (s *StrictControllerServer) userIDFromContext(ctx context.Context) (shared.UserID, bool) {
	if ginCtx, ok := ginContext(ctx); ok {
		if value, exists := ginCtx.Get("user_id"); exists {
			if userID, ok := value.(shared.UserID); ok {
				return userID, true
			}
		}
	}
	return "", false
}

func ginContext(ctx context.Context) (*gin.Context, bool) {
	if g, ok := ctx.(*gin.Context); ok {
		return g, true
	}
	return nil, false
}

func toAPIGrumble(resp *controller.GrumbleResponse) Grumble {
	var hasVibed *bool
	if resp.HasVibed != nil {
		value := *resp.HasVibed
		hasVibed = &value
	}

	return Grumble{
		GrumbleID:         openapi_types.UUID(resp.GrumbleID),
		UserID:            openapi_types.UUID(resp.UserID),
		Content:           resp.Content,
		ToxicLevel:        resp.ToxicLevel,
		VibeCount:         resp.VibeCount,
		PurifiedThreshold: resp.PurifiedThreshold,
		IsPurified:        resp.IsPurified,
		PostedAt:          resp.PostedAt,
		ExpiresAt:         resp.ExpiresAt,
		IsEventGrumble:    resp.IsEventGrumble,
		HasVibed:          hasVibed,
	}
}

func toAPIVibe(resp *controller.AddVibeResponse) Vibe {
	return Vibe{
		VibeID:    resp.VibeID,
		GrumbleID: openapi_types.UUID(resp.GrumbleID),
		UserID:    openapi_types.UUID(resp.UserID),
		VibeType:  VibeVibeType(resp.VibeType),
		VotedAt:   resp.VotedAt,
	}
}

func toAPIAnonymousUser(resp *controller.MyProfileResponse) AnonymousUser {
	anon := AnonymousUser{
		UserID:       openapi_types.UUID(resp.UserID),
		VirtuePoints: resp.VirtuePoints,
		CreatedAt:    resp.CreatedAt,
	}
	if resp.ProfileTitle != nil {
		anon.ProfileTitle = nullable.NewNullableWithValue(*resp.ProfileTitle)
	}
	if resp.VirtueRank != "" {
		rank := AnonymousUserVirtueRank(resp.VirtueRank)
		anon.VirtueRank = &rank
	}
	return anon
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "invalid") || strings.Contains(msg, "must be") || strings.Contains(msg, "between")
}

func errorResponse(code, message string) ErrorResponse {
	return ErrorResponse{Error: code, Message: message}
}
