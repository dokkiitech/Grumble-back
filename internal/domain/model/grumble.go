package model

import (
	"time"

	"github.com/takagiyuuki/grumble-back/internal/domain/value"
)

// Grumble は愚痴投稿を表す集約ルート。
type Grumble struct {
	ID             value.GrumbleID  `db:"grumble_id"`
	UserID         value.UserID     `db:"user_id"`
	Content        string           `db:"content"`
	ToxicLevel     value.ToxicLevel `db:"toxic_level"`
	VibeCount      int              `db:"vibe_count"`
	IsPurified     bool             `db:"is_purified"`
	PostedAt       time.Time        `db:"posted_at"`
	ExpiresAt      time.Time        `db:"expires_at"`
	IsEventGrumble bool             `db:"is_event_grumble"`
}
