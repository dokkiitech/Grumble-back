//go:build ignore

package model

import (
	"time"
)

// EventType はイベントの種類。
type EventType string

const (
	// EventTypeDaionryo は大怨霊討伐イベントを表す。
	EventTypeDaionryo EventType = "DAIONRYO"
	// EventTypeOtakinage はお焚き上げイベントを表す。
	EventTypeOtakinage EventType = "OTAKINAGE"
)

// Event は期間限定イベントの状態を保持する。
type Event struct {
	ID        value.EventID `db:"event_id"`
	Name      string        `db:"event_name"`
	Type      EventType     `db:"event_type"`
	StartTime time.Time     `db:"start_time"`
	EndTime   time.Time     `db:"end_time"`
	CurrentHP int           `db:"current_hp"`
	MaxHP     int           `db:"max_hp"`
	IsActive  bool          `db:"is_active"`
}
