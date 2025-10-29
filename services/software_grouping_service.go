package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type SoftwareLicense struct {
	Name              string  `json:"name"`
	THBPricePerYear   float64 `json:"thb_price_per_year"`
	UsedByCompanyName string  `json:"used_by_company_name"`
}

type GroupedSoftwareResponse struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Items          []map[string]any `json:"items"`
	CommonFeatures []string         `json:"common_features"`
}

type SoftwareGroupingService struct {
	client *openai.Client
	model  string
}

const (
	groupingModelDefault = "gpt-4o-mini"
	maxTokensDefault     = 4000
	requestTimeout       = 60 * time.Second
)

func NewSoftwareGroupingService(apiKey string) (*SoftwareGroupingService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("missing OpenAI API key")
	}

	client := openai.NewClient(apiKey)

	return &SoftwareGroupingService{
		client: client,
		model:  groupingModelDefault,
	}, nil
}

func (s *SoftwareGroupingService) GroupSoftwareLicenses(ctx context.Context, licenses []SoftwareLicense) ([]GroupedSoftwareResponse, error) {
	// Use compact JSON to reduce token count and speed up processing
	inputJSON, err := json.Marshal(licenses)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	systemPrompt := `You are an expert enterprise software analyst. Normalize software names (fix typos, unify variants), group them into logical enterprise categories (Design Tools, Developer Tools, Project Management, Business Intelligence, Communication & Collaboration, IT Service Management, Cloud Platforms, Productivity Suites), preserve company usage as separate items with their own used_by_company_name and thb_price_per_year, and provide real logo/image URLs from public sources (Wikipedia, Wikimedia, official sites).

Return ONLY valid JSON array (no markdown/code blocks):
[{"name":"Group name","description":"Brief description","items":[{"name":"Normalized name","image":"real_url","thb_price_per_year":number,"used_by_company_name":"Company"}],"common_features":["Feature1","Feature2","Feature3"]}]

Rules: Valid JSON only. Normalize identical software names consistently. Keep separate entries per subsidiary. Use real image URLs. Include ≥3 common features per group.`

	userPrompt := fmt.Sprintf("Group these licenses:\n%s", string(inputJSON))

	// Add timeout to context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: 0.1, // Lower temperature for faster, more deterministic output
		MaxTokens:   maxTokensDefault,
	})

	if err != nil {
		// Check if it's a timeout error
		if ctxWithTimeout.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("openai request timeout after %v: %w", requestTimeout, err)
		}
		return nil, fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no valid response from OpenAI")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Clean up JSON if wrapped in markdown code blocks
	if strings.HasPrefix(content, "```json") || strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	var result []GroupedSoftwareResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output (content length: %d): %w", len(content), err)
	}

	return result, nil
}
