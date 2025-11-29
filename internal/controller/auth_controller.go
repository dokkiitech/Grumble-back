package controller

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/google/uuid"
)

// AuthController handles authentication-related use cases.
type AuthController struct {
	authAnonymousUC     *usecase.AuthAnonymousUseCase
	userQueryUC         *usecase.UserQueryUseCase
	logger              logging.Logger
	rankingLimitDefault int
	rankingLimitMin     int
	rankingLimitMax     int
}

// NewAuthController creates a new AuthController.
func NewAuthController(
	authAnonymousUC *usecase.AuthAnonymousUseCase,
	userQueryUC *usecase.UserQueryUseCase,
	logger logging.Logger,
	rankingLimitDefault int,
	rankingLimitMin int,
	rankingLimitMax int,
) *AuthController {
	if rankingLimitMax < rankingLimitMin {
		rankingLimitMax = rankingLimitMin
	}
	if rankingLimitDefault < rankingLimitMin || rankingLimitDefault > rankingLimitMax {
		rankingLimitDefault = rankingLimitMin
	}
	return &AuthController{
		authAnonymousUC:     authAnonymousUC,
		userQueryUC:         userQueryUC,
		logger:              logger,
		rankingLimitDefault: rankingLimitDefault,
		rankingLimitMin:     rankingLimitMin,
		rankingLimitMax:     rankingLimitMax,
	}
}

// MyProfileResponse represents the response for GET /users/me.
type MyProfileResponse struct {
	UserID       uuid.UUID
	VirtuePoints int
	VirtueRank   string
	CreatedAt    time.Time
	ProfileTitle *string
}

// GetMyProfile fetches the authenticated user's profile.
func (ctrl *AuthController) GetMyProfile(ctx context.Context, userID shared.UserID) (*MyProfileResponse, error) {
	user, err := ctrl.userQueryUC.GetMyProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(string(user.UserID))
	if err != nil {
		ctrl.logger.ErrorContext(ctx, "Failed to parse user UUID", "error", err)
		return nil, err
	}

	return &MyProfileResponse{
		UserID:       userUUID,
		VirtuePoints: user.VirtuePoints,
		VirtueRank:   string(user.Rank()),
		CreatedAt:    user.CreatedAt,
		ProfileTitle: user.ProfileTitle,
	}, nil
}

// ResolveRankingLimit normalises the ranking limit respecting configured bounds.
func (ctrl *AuthController) ResolveRankingLimit(requested *int) int {
	if requested == nil {
		return ctrl.rankingLimitDefault
	}
	value := *requested
	if value < ctrl.rankingLimitMin {
		return ctrl.rankingLimitMin
	}
	if value > ctrl.rankingLimitMax {
		return ctrl.rankingLimitMax
	}
	return value
}

// RankingResponse represents an entry in the bodhisattva ranking list.
type RankingResponse struct {
	UserID       uuid.UUID
	VirtuePoints int
	VirtueRank   string
	CreatedAt    time.Time
	ProfileTitle *string
}

// GetBodhisattvaRankings fetches the ranking list.
func (ctrl *AuthController) GetBodhisattvaRankings(ctx context.Context, limit int) ([]*RankingResponse, error) {
	if limit <= 0 {
		limit = ctrl.rankingLimitDefault
	}

	users, err := ctrl.userQueryUC.GetBodhisattvaRankings(ctx, limit)
	if err != nil {
		return nil, err
	}

	rankings := make([]*RankingResponse, len(users))
	for i, u := range users {
		userUUID, err := uuid.Parse(string(u.UserID))
		if err != nil {
			ctrl.logger.ErrorContext(ctx, "Failed to parse user UUID", "error", err)
			return nil, err
		}
		rankings[i] = &RankingResponse{
			UserID:       userUUID,
			VirtuePoints: u.VirtuePoints,
			VirtueRank:   string(u.Rank()),
			CreatedAt:    u.CreatedAt,
			ProfileTitle: u.ProfileTitle,
		}
	}

	return rankings, nil
}
