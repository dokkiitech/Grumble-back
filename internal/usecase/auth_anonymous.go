package usecase

import (
	"context"
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/domain/user"
	"github.com/google/uuid"
)

// AuthAnonymousUseCase handles anonymous user authentication
type AuthAnonymousUseCase struct {
	userRepo user.Repository
}

// NewAuthAnonymousUseCase creates a new AuthAnonymousUseCase
func NewAuthAnonymousUseCase(userRepo user.Repository) *AuthAnonymousUseCase {
	return &AuthAnonymousUseCase{
		userRepo: userRepo,
	}
}

// AuthAnonymousRequest represents the input for anonymous authentication
type AuthAnonymousRequest struct {
	UserID shared.UserID // Device-generated UUID
}

// AuthAnonymousResponse represents the authentication result
type AuthAnonymousResponse struct {
	User      *user.AnonymousUser
	IsNewUser bool
}

// Authenticate authenticates or creates an anonymous user
func (uc *AuthAnonymousUseCase) Authenticate(ctx context.Context, req AuthAnonymousRequest) (*AuthAnonymousResponse, error) {
	// Validate UUID format
	if _, err := uuid.Parse(string(req.UserID)); err != nil {
		return nil, &shared.ValidationError{
			Field:   "user_id",
			Message: "user_id must be a valid UUID",
		}
	}

	// Try to find existing user
	existingUser, err := uc.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		// if not NotFound, bubble up
		if _, ok := err.(*shared.NotFoundError); !ok {
			return nil, err
		}
		// Not found: create new user
	} else {
		// User found, return existing user
		return &AuthAnonymousResponse{
			User:      existingUser,
			IsNewUser: false,
		}, nil
	}

	// Create new anonymous user
	newUser := &user.AnonymousUser{
		UserID:       req.UserID,
		VirtuePoints: 0,
		CreatedAt:    time.Now(),
		ProfileTitle: nil,
	}

	// Validate before creating
	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	// Persist to repository
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return &AuthAnonymousResponse{
		User:      newUser,
		IsNewUser: true,
	}, nil
}
