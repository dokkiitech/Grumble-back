package grumble

// VibeRank represents a label derived from vibe_count.
type VibeRank string

const (
	VibeRankNone       VibeRank = ""
	VibeRankMinarai    VibeRank = "見習い行者"
	VibeRankMiroku     VibeRank = "弥勒"
	VibeRankJizo       VibeRank = "地蔵"
	VibeRankDaibosatsu VibeRank = "大菩薩"
)

// RankFromVibeCount returns the rank label determined by the vibe count.
func RankFromVibeCount(vibeCount int) VibeRank {
	switch {
	case vibeCount >= 40:
		return VibeRankDaibosatsu
	case vibeCount >= 30:
		return VibeRankJizo
	case vibeCount >= 20:
		return VibeRankMiroku
	case vibeCount >= 10:
		return VibeRankMinarai
	default:
		return VibeRankNone
	}
}
