package service

import "github.com/dokkiitech/grumble-back/internal/domain/user"

// VirtueService contains rules for virtue point calculations.
type VirtueService struct{}

// NewVirtueService instantiates a VirtueService.
func NewVirtueService() *VirtueService {
	return &VirtueService{}
}

// CalculateVirtuePoints returns the virtue points total after giving a vibe.
func (s *VirtueService) CalculateVirtuePoints(u *user.AnonymousUser, vibeGiven bool) int {
	if u == nil {
		return 0
	}
	if vibeGiven {
		return u.VirtuePoints + 1
	}
	return u.VirtuePoints
}
