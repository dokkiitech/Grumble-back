package controller

import (
	"github.com/dokkiitech/grumble-back/internal/usecase"
)

// TimelineResponse represents a timeline response
type TimelineResponse struct {
	Grumbles []*GrumbleResponse `json:"grumbles"`
	Total    int                `json:"total"`
}

// TimelinePresenter converts timeline results to API responses
type TimelinePresenter struct {
	grumblePresenter *GrumblePresenter
}

// NewTimelinePresenter creates a new TimelinePresenter
func NewTimelinePresenter(grumblePresenter *GrumblePresenter) *TimelinePresenter {
	return &TimelinePresenter{
		grumblePresenter: grumblePresenter,
	}
}

// ToAPITimelineResponse converts usecase timeline response to API format
func (p *TimelinePresenter) ToAPITimelineResponse(resp *usecase.TimelineResponse) (*TimelineResponse, error) {
	grumbles, err := p.grumblePresenter.ToAPIGrumbles(resp.Grumbles)
	if err != nil {
		return nil, err
	}

	return &TimelineResponse{
		Grumbles: grumbles,
		Total:    resp.TotalCount,
	}, nil
}
