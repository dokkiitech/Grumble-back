package value

import "strconv"

// 各テーブルで共有するID型を定義する。
type (
	GrumbleID  string
	UserID     string
	VibeID     int64
	PollID     int64
	PollVoteID int64
	EventID    int64
)

// String は内部表現をそのまま返すヘルパ。
func (id GrumbleID) String() string { return string(id) }

// String は内部表現を文字列化して返すヘルパ。
func (id VibeID) String() string { return strconv.FormatInt(int64(id), 10) }

// Int64 は内部表現を int64 で返す。
func (id VibeID) Int64() int64 { return int64(id) }

// String は内部表現をそのまま返すヘルパ。
func (id UserID) String() string { return string(id) }

// String は内部表現を文字列化して返すヘルパ。
func (id PollID) String() string { return strconv.FormatInt(int64(id), 10) }

// Int64 は内部表現を int64 で返す。
func (id PollID) Int64() int64 { return int64(id) }

// String は内部表現を文字列化して返すヘルパ。
func (id PollVoteID) String() string { return strconv.FormatInt(int64(id), 10) }

// Int64 は内部表現を int64 で返す。
func (id PollVoteID) Int64() int64 { return int64(id) }

// String は内部表現を文字列化して返すヘルパ。
func (id EventID) String() string { return strconv.FormatInt(int64(id), 10) }

// Int64 は内部表現を int64 で返す。
func (id EventID) Int64() int64 { return int64(id) }
