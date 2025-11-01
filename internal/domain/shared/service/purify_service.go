package service

import (
	"errors"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
)

// PurifyService encapsulates purification rules for grumbles.
type PurifyService struct {
	threshold int
}

// NewPurifyService creates a new PurifyService with a purification threshold.
func NewPurifyService(threshold int) *PurifyService {
	return &PurifyService{threshold: threshold}
}

// Threshold returns the configured purification threshold.
func (s *PurifyService) Threshold() int {
	return s.threshold
}

// ShouldPurify determines if the grumble meets purification criteria.
func (s *PurifyService) ShouldPurify(g *grumble.Grumble) bool {
	if g == nil {
		return false
	}
	return g.CanBePurified(s.threshold)
}

// Purify marks the grumble as purified if it satisfies the rules.
func (s *PurifyService) Purify(g *grumble.Grumble) error {
	if !s.ShouldPurify(g) {
		return errors.New("grumble does not meet purification criteria")
	}
	g.Purify()
	return nil
}
