package model

import (
	"time"
)

// VibeType は共感の種類を表す。
type VibeType string

const (
	// VibeTypeWakaru はデフォルトの「わかる…」リアクション。
	VibeTypeWakaru VibeType = "WAKARU"
)

// Vibe は投稿に対する共感履歴。
type Vibe struct {
	ID        value.VibeID    `db:"vibe_id"`
	GrumbleID value.GrumbleID `db:"grumble_id"`
	UserID    value.UserID    `db:"user_id"`
	Type      VibeType        `db:"vibe_type"`
	VotedAt   time.Time       `db:"voted_at"`
}
