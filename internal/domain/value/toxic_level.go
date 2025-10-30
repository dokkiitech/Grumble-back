package value

import "fmt"

// ToxicLevel は投稿時に自己申告する毒レベルを表す。
// 1〜5 の範囲であることを保証するため、Validate を併用する。
type ToxicLevel int

const (
	// ToxicLevelMin は許容される最小値。
	ToxicLevelMin ToxicLevel = 1
	// ToxicLevelMax は許容される最大値。
	ToxicLevelMax ToxicLevel = 5
)

// Validate は毒レベルが定義済み範囲内であることを確認する。
func (t ToxicLevel) Validate() error {
	if t < ToxicLevelMin || t > ToxicLevelMax {
		return fmt.Errorf("toxic level must be between %d and %d: %d", ToxicLevelMin, ToxicLevelMax, t)
	}
	return nil
}

// Int は整数値を返すヘルパ。
func (t ToxicLevel) Int() int {
	return int(t)
}
