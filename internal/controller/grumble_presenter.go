package controller

import (
	"time"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/google/uuid"
)

// GrumbleResponse represents a grumble in API responses
type GrumbleResponse struct {
	GrumbleID         uuid.UUID `json:"grumble_id"`
	UserID            uuid.UUID `json:"user_id"`
	Content           string    `json:"content"`
	ToxicLevel        int       `json:"toxic_level"`
	VibeCount         int       `json:"vibe_count"`
	PurifiedThreshold int       `json:"purified_threshold"`
	IsPurified        bool      `json:"is_purified"`
	PostedAt          time.Time `json:"posted_at"`
	ExpiresAt         time.Time `json:"expires_at"`
	IsEventGrumble    bool      `json:"is_event_grumble"`
	HasVibed          *bool     `json:"has_vibed,omitempty"`
}

// GrumblePresenter converts domain grumbles to API responses
type GrumblePresenter struct{}

// NewGrumblePresenter creates a new GrumblePresenter
func NewGrumblePresenter() *GrumblePresenter {
	return &GrumblePresenter{}
}

// ToAPIGrumble converts a domain Grumble to API Grumble response
func (p *GrumblePresenter) ToAPIGrumble(g *grumble.Grumble) (*GrumbleResponse, error) {
	grumbleUUID, err := uuid.Parse(string(g.GrumbleID))
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.Parse(string(g.UserID))
	if err != nil {
		return nil, err
	}

	return &GrumbleResponse{
		GrumbleID:         grumbleUUID,
		UserID:            userUUID,
		Content:           g.Content,
		ToxicLevel:        int(g.ToxicLevel),
		VibeCount:         g.VibeCount,
		PurifiedThreshold: g.PurifiedThreshold,
		IsPurified:        g.IsPurified,
		PostedAt:          g.PostedAt,
		ExpiresAt:         g.ExpiresAt,
		IsEventGrumble:    g.IsEventGrumble,
		HasVibed:          g.HasVibed,
	}, nil
}

// ToAPIGrumbles converts multiple domain Grumbles to API response array
func (p *GrumblePresenter) ToAPIGrumbles(grumbles []*grumble.Grumble) ([]*GrumbleResponse, error) {
	result := make([]*GrumbleResponse, len(grumbles))
	for i, g := range grumbles {
		apiGrumble, err := p.ToAPIGrumble(g)
		if err != nil {
			return nil, err
		}
		result[i] = apiGrumble
	}
	return result, nil
}
