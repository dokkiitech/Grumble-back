package usecase

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/domain/user"
)

// UserQueryUseCase handles read-only user queries for controllers
type UserQueryUseCase struct {
	userRepo user.Repository
}

func NewUserQueryUseCase(userRepo user.Repository) *UserQueryUseCase {
	return &UserQueryUseCase{userRepo: userRepo}
}

// GetMyProfile retrieves the profile of the given user ID
func (uc *UserQueryUseCase) GetMyProfile(ctx context.Context, id shared.UserID) (*user.AnonymousUser, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// GetBodhisattvaRankings retrieves top users by virtue points
func (uc *UserQueryUseCase) GetBodhisattvaRankings(ctx context.Context, limit int) ([]*user.AnonymousUser, error) {
	return uc.userRepo.FindTopByVirtuePoints(ctx, limit)
}
