package user

// VirtueRank represents a label derived from virtue points (number of vibes given).
type VirtueRank string

const (
	VirtueRankNone       VirtueRank = ""
	VirtueRankMinarai    VirtueRank = "見習い行者"
	VirtueRankMiroku     VirtueRank = "弥勒"
	VirtueRankJizo       VirtueRank = "地蔵"
	VirtueRankDaibosatsu VirtueRank = "大菩薩"
)

// RankFromVirtuePoints returns the rank label based on virtue points.
// Thresholds are 10, 20, 30, 40.
func RankFromVirtuePoints(points int) VirtueRank {
	switch {
	case points >= 40:
		return VirtueRankDaibosatsu
	case points >= 30:
		return VirtueRankJizo
	case points >= 20:
		return VirtueRankMiroku
	case points >= 10:
		return VirtueRankMinarai
	default:
		return VirtueRankNone
	}
}
