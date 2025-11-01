package shared

// Strongly-typed IDs
type UserID string
type GrumbleID string
type VibeID int
type EventID int

// ToxicLevel represents the self-reported toxicity level (1-5)
type ToxicLevel int

const (
	ToxicLevel1 ToxicLevel = 1 // Mild annoyance
	ToxicLevel2 ToxicLevel = 2 // Moderate frustration
	ToxicLevel3 ToxicLevel = 3 // Significant anger
	ToxicLevel4 ToxicLevel = 4 // Major outrage
	ToxicLevel5 ToxicLevel = 5 // Extreme fury
)

func (t ToxicLevel) Validate() error {
	if t < 1 || t > 5 {
		return &ValidationError{Field: "toxic_level", Message: "must be 1-5"}
	}
	return nil
}

// VibeType represents the type of empathy reaction
type VibeType string

const (
	VibeTypeWakaru VibeType = "WAKARU" // Primary empathy type
)

func (v VibeType) Validate() error {
	if v != VibeTypeWakaru {
		return &ValidationError{Field: "vibe_type", Message: "only WAKARU supported"}
	}
	return nil
}
