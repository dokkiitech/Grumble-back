package grumble

import "context"

// ContentModerationPrompt is the prompt template for Gemini API to moderate content
const ContentModerationPrompt = `以下の投稿内容を審査し、JSON形式で判定結果を出力して。
**jsonのみ出力してください**

# 不適切と判定する基準
1. 誹謗中傷・攻撃的な表現
2. 差別的な内容（人種、性別、宗教、国籍、障害等）
3. 個人情報（氏名、住所、電話番号、メールアドレス等）
4. 違法行為の助長、自傷行為の勧誘

# 出力形式
{
  "is_appropriate": true/false,
  "reason": "理由"
}

# 投稿内容
%s`

// ModerationResult represents the result of content moderation
type ModerationResult struct {
	IsAppropriate bool   `json:"is_appropriate"`
	Reason        string `json:"reason"`
}

// ContentFilterClient is an interface for content moderation
type ContentFilterClient interface {
	// FilterContent checks if the content is appropriate
	// Returns ModerationResult with the filtering decision
	FilterContent(ctx context.Context, content string) (*ModerationResult, error)
}
