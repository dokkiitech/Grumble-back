package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dokkiitech/grumble-back/internal/domain/grumble"
	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"google.golang.org/genai"
)

// GeminiClient implements grumble.ContentFilterClient using Gemini API
type GeminiClient struct {
	apiKey string
	model  string
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(apiKey, model string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		model:  model,
	}
}

// FilterContent implements grumble.ContentFilterClient
func (c *GeminiClient) FilterContent(ctx context.Context, content string) (*grumble.ModerationResult, error) {
	// If API key is not provided, skip external moderation to avoid hard failure in dev/local.
	if strings.TrimSpace(c.apiKey) == "" {
		return &grumble.ModerationResult{
			IsAppropriate: true,
			Reason:        "gemini_api_key_missing_skipped",
		}, nil
	}

	// Create client with API key from environment
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: c.apiKey,
	})
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to create Gemini client",
			Err:     err,
		}
	}

	// Build prompt
	prompt := fmt.Sprintf(grumble.ContentModerationPrompt, content)

	// Generate content with JSON response format
	result, err := client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, &shared.InternalError{
			Message: "failed to generate content from Gemini",
			Err:     err,
		}
	}

	// Get response text
	responseText := result.Text()
	if responseText == "" {
		return nil, &shared.InternalError{
			Message: "empty response from Gemini",
		}
	}

	// Clean up response text (remove markdown code blocks if present)
	responseText = strings.TrimSpace(responseText)
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// Unmarshal and validate
	var moderationResult grumble.ModerationResult
	if err := json.Unmarshal([]byte(responseText), &moderationResult); err != nil {
		return nil, &shared.InternalError{
			Message: fmt.Sprintf("failed to parse Gemini response: %s", responseText),
			Err:     err,
		}
	}

	// Validate result
	if moderationResult.Reason == "" {
		return nil, &shared.InternalError{
			Message: fmt.Sprintf("invalid moderation result: reason is empty. Response: %s", responseText),
		}
	}

	return &moderationResult, nil
}
