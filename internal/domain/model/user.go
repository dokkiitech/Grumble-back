package model

import (
	"time"

	"github.com/takagiyuuki/grumble-back/internal/domain/value"
)

// AnonymousUser は匿名利用者アカウント。
type AnonymousUser struct {
	ID           value.UserID `db:"user_id"`
	VirtuePoints int          `db:"virtue_points"`
	CreatedAt    time.Time    `db:"created_at"`
	ProfileTitle *string      `db:"profile_title"`
}
