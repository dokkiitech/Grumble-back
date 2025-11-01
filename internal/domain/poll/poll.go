//go:build ignore

package model

// Poll は投稿に紐づく二択投票（将来機能）。
type Poll struct {
	ID        value.PollID    `db:"poll_id"`
	GrumbleID value.GrumbleID `db:"grumble_id"`
	Question  string          `db:"question"`
	Option1   string          `db:"option_1"`
	Option2   string          `db:"option_2"`
}

// PollVote は投票結果（将来機能）。
type PollVote struct {
	ID             value.PollVoteID `db:"poll_vote_id"`
	PollID         value.PollID     `db:"poll_id"`
	UserID         value.UserID     `db:"user_id"`
	SelectedOption int              `db:"selected_option"`
}
