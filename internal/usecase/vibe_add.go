package usecase

import (
	"context"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	sharedservice "github.com/dokkiitech/grumble-back/internal/domain/shared/service"
	"github.com/dokkiitech/grumble-back/internal/domain/user"
	"github.com/dokkiitech/grumble-back/internal/domain/vibe"
)

// VibeAddUseCase handles giving vibes to grumbles.
type VibeAddUseCase struct {
	grumbleRepo grumble.Repository
	vibeRepo    vibe.Repository
	userRepo    user.Repository
	purifySvc   *sharedservice.PurifyService
	virtueSvc   *sharedservice.VirtueService
}

// NewVibeAddUseCase constructs a VibeAddUseCase.
func NewVibeAddUseCase(
	grumbleRepo grumble.Repository,
	vibeRepo vibe.Repository,
	userRepo user.Repository,
	purifySvc *sharedservice.PurifyService,
	virtueSvc *sharedservice.VirtueService,
) *VibeAddUseCase {
	return &VibeAddUseCase{
		grumbleRepo: grumbleRepo,
		vibeRepo:    vibeRepo,
		userRepo:    userRepo,
		purifySvc:   purifySvc,
		virtueSvc:   virtueSvc,
	}
}

// AddVibeRequest represents inputs for giving a vibe.
type AddVibeRequest struct {
	GrumbleID shared.GrumbleID
	UserID    shared.UserID
	VibeType  shared.VibeType
}

// AddVibeResponse represents the outcome of giving a vibe.
type AddVibeResponse struct {
	Vibe         *vibe.Vibe
	VibeCount    int
	VirtuePoints int
	IsPurified   bool
}

// Add processes a vibe action.
func (uc *VibeAddUseCase) Add(ctx context.Context, req AddVibeRequest) (*AddVibeResponse, error) {
	vType := req.VibeType
	if vType == "" {
		vType = shared.VibeTypeWakaru
	}

	// Fetch grumble to ensure it exists and get current state.
	grumbleEntity, err := uc.grumbleRepo.FindByID(ctx, req.GrumbleID)
	if err != nil {
		return nil, err
	}

	if grumbleEntity.IsPurified {
		return nil, &shared.ValidationError{
			Field:   "grumble",
			Message: "grumble already purified",
		}
	}

	if grumbleEntity.UserID == req.UserID {
		return nil, &shared.ValidationError{
			Field:   "grumble_id",
			Message: "cannot vibe your own grumble",
		}
	}

	// Prevent duplicate vibes
	exists, err := uc.vibeRepo.Exists(ctx, req.GrumbleID, req.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &shared.DuplicateVibeError{
			GrumbleID: string(req.GrumbleID),
			UserID:    string(req.UserID),
		}
	}

	// Load user for virtue calculation
	userEntity, err := uc.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	expectedVirtuePoints := uc.virtueSvc.CalculateVirtuePoints(userEntity, true)

	v := &vibe.Vibe{
		GrumbleID: req.GrumbleID,
		UserID:    req.UserID,
		Type:      vType,
	}

	if err := v.Validate(); err != nil {
		return nil, err
	}

	createResult, err := uc.vibeRepo.Create(ctx, v)
	if err != nil {
		return nil, err
	}

	// Update in-memory representations with new state.
	grumbleEntity.VibeCount = createResult.VibeCount
	if createResult.VirtuePoints > 0 {
		userEntity.VirtuePoints = createResult.VirtuePoints
	} else {
		userEntity.VirtuePoints = expectedVirtuePoints
	}

	// Check purification threshold.
	if uc.purifySvc != nil && uc.purifySvc.ShouldPurify(grumbleEntity) {
		if err := uc.purifySvc.Purify(grumbleEntity); err != nil {
			return nil, err
		}
		if err := uc.grumbleRepo.Update(ctx, grumbleEntity); err != nil {
			return nil, err
		}
	}

	return &AddVibeResponse{
		Vibe:         createResult.Vibe,
		VibeCount:    grumbleEntity.VibeCount,
		VirtuePoints: userEntity.VirtuePoints,
		IsPurified:   grumbleEntity.IsPurified,
	}, nil
}
