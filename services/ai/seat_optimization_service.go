package ai

import (
	dbRepo "asset-pulse-api/repositories/database"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// SeatOptimizationService handles AI-powered seat optimization recommendations
type SeatOptimizationService struct {
	client      *openai.Client
	db          *gorm.DB
	dbRepo      dbRepo.DatabaseRepository
	initialized bool
}

// OptimizationAnalysisRequest for analyzing seat optimization opportunities
type OptimizationAnalysisRequest struct {
	CompanyCode    string
	DepartmentCode string
	AppName        string
	Limit          int
	Offset         int
}

// OptimizationAnalysisResponse contains AI-generated recommendations
type OptimizationAnalysisResponse struct {
	Recommendations []AIRecommendation
	Confidence      float64
}

// AIRecommendation represents an AI-generated optimization recommendation
type AIRecommendation struct {
	Action           string  `json:"action"`
	AppName          string  `json:"app_name"`
	Department       string  `json:"department"`
	Rationale        string  `json:"rationale"`
	Priority         int     `json:"priority"`
	RiskLevel        string  `json:"risk_level"`
	PotentialSavings float64 `json:"potential_savings"`
	FromDepartment   string  `json:"from_department"`
	ToDepartment     string  `json:"to_department"`
	DowngradeFrom    string  `json:"downgrade_from"`
	DowngradeTo      string  `json:"downgrade_to"`
}

const (
	analyzeModel = "gpt-4o-mini"
)

// NewSeatOptimizationService creates a new seat optimization service
func NewSeatOptimizationService(apiKey string, db *gorm.DB, dbRepo dbRepo.DatabaseRepository) *SeatOptimizationService {
	client := openai.NewClient(apiKey)
	return &SeatOptimizationService{
		client:      client,
		db:          db,
		dbRepo:      dbRepo,
		initialized: false,
	}
}

// AnalyzeOpportunities analyzes license usage data and generates AI recommendations
func (s *SeatOptimizationService) AnalyzeOpportunities(ctx context.Context, req *OptimizationAnalysisRequest) (*OptimizationAnalysisResponse, error) {
	// Fetch usage analytics
	usageData, err := s.dbRepo.GetLicenseUsageAnalytics(ctx, req.CompanyCode, req.DepartmentCode, req.AppName)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage analytics: %w", err)
	}

	if len(usageData) == 0 {
		return &OptimizationAnalysisResponse{
			Recommendations: []AIRecommendation{},
			Confidence:      0,
		}, nil
	}

	// Fetch inactive users without limit (use pagination in batches)
	allInactiveUsers, err := s.dbRepo.GetInactiveUsers(ctx, req.CompanyCode, 0, 90, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive users: %w", err)
	}

	// Convert usage data to JSON
	usageJSON, err := json.Marshal(usageData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal usage data: %w", err)
	}

	// Process in batches to prevent huge AI prompts
	batchSize := 500
	allRecommendations := []AIRecommendation{}

	for i := 0; i < len(allInactiveUsers); i += batchSize {
		end := i + batchSize
		if end > len(allInactiveUsers) {
			end = len(allInactiveUsers)
		}

		batch := allInactiveUsers[i:end]

		// Marshal current batch
		batchJSON, err := json.Marshal(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal inactive users batch: %w", err)
		}

		// Create AI prompt
		prompt := fmt.Sprintf(`You are a license optimization AI assistant. Analyze the following data and recommend license optimization opportunities.

License Usage Data:
%s

Inactive Users (not used in 90+ days) - Batch %d/%d:
%s

For each recommendation, provide:
1. Action type: "revoke", "reallocate", or "downgrade"
2. Rationale: Why this action is recommended (be specific and concise)
3. Priority: 1-3 (3 = high priority)
4. Risk Level: "Low", "Medium", or "High"
5. For reallocate: specify from_department and to_department (use department_code)
6. For downgrade: specify downgrade_from and downgrade_to
7. Estimated potential savings per month

IMPORTANT: Use department_code (like "ENG", "MKT") for department field, not department name.

Focus on:
- License assignments with high inactive user counts
- Opportunities to reallocate from inactive to requesting departments
- Tier downgrades for users with minimal feature usage
- Prioritize high-value savings with low risk

Return ONLY valid JSON format:
{
  "recommendations": [
    {
      "action": "revoke",
      "app_name": "...",
      "department": "ENG",  // Use department_code
      "rationale": "...",
      "priority": 2,
      "risk_level": "Low",
      "potential_savings": 2400.0,
      "from_department": "",
      "to_department": "",
      "downgrade_from": "",
      "downgrade_to": ""
    }
  ],
  "confidence": 0.85
}`, string(usageJSON), i/batchSize+1, (len(allInactiveUsers)+batchSize-1)/batchSize, string(batchJSON))

		// Call OpenAI API
		resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       analyzeModel,
			Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
			Temperature: 0.3,
		})

		if err != nil {
			return nil, fmt.Errorf("openai api error: %w", err)
		}

		if len(resp.Choices) == 0 {
			continue // Skip this batch if no response
		}

		content := strings.TrimSpace(resp.Choices[0].Message.Content)

		// Clean up content est (remove markdown code blocks if present)
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)

		// Parse JSON response
		var batchResponse struct {
			Recommendations []AIRecommendation `json:"recommendations"`
			Confidence      float64            `json:"confidence"`
		}

		if err := json.Unmarshal([]byte(content), &batchResponse); err != nil {
			continue // Skip this batch if parse fails
		}

		// Append to all recommendations
		allRecommendations = append(allRecommendations, batchResponse.Recommendations...)
	}

	// Limit results based on request
	if req.Limit > 0 && len(allRecommendations) > req.Limit {
		allRecommendations = allRecommendations[:req.Limit]
	}

	return &OptimizationAnalysisResponse{
		Recommendations: allRecommendations,
		Confidence:      0.85, // Average confidence
	}, nil
}
