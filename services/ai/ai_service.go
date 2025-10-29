package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// AI Service Interface
type AIService interface {
	GenerateJDRecommendations(ctx context.Context, req *JDRecommendationRequest) (*JDRecommendationResponse, error)
	GenerateConsolidationMemo(ctx context.Context, req *ConsolidationMemoRequest) (*ConsolidationMemoResponse, error)
	CalculateSoftwareSimilarity(ctx context.Context, req *SimilarityRequest) (*SimilarityResponse, error)
	GetAppRecommendationScores(ctx context.Context, prompt string) (string, error)
	GenerateShortDescription(ctx context.Context, appName, category string, score float64) (string, error)
}

// Mock AI Service Implementation
type MockAIService struct {
	// In production, this would contain OpenAI/Anthropic client
}

func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

// JD Recommendation Request/Response
type JDRecommendationRequest struct {
	JobTitle       string   `json:"job_title"`
	JobDescription string   `json:"job_description"`
	Department     string   `json:"department"`
	CompanyCode    string   `json:"company_code"`
	Experience     string   `json:"experience"`
	Skills         []string `json:"skills"`
}

type JDRecommendationResponse struct {
	Recommendations []SoftwareRecommendation `json:"recommendations"`
	Confidence      float64                  `json:"confidence"`
	ProcessingTime  time.Duration            `json:"processing_time"`
}

type SoftwareRecommendation struct {
	AppName        string   `json:"app_name"`
	Category       string   `json:"category"`
	Tier           string   `json:"tier"`
	RelevanceScore float64  `json:"relevance_score"`
	Cost           float64  `json:"cost"`
	Rationale      string   `json:"rationale"`
	Features       []string `json:"features"`
	Alternatives   []string `json:"alternatives"`
}

// Consolidation Memo Request/Response
type ConsolidationMemoRequest struct {
	Apps          []string `json:"apps"`
	Subsidiaries  []string `json:"subsidiaries"`
	CurrentCost   float64  `json:"current_cost"`
	TargetSavings float64  `json:"target_savings"`
	Context       string   `json:"context"`
}

type ConsolidationMemoResponse struct {
	Memo           string  `json:"memo"`
	Savings        float64 `json:"savings"`
	RiskAssessment string  `json:"risk_assessment"`
	Timeline       string  `json:"timeline"`
	Confidence     float64 `json:"confidence"`
}

// Similarity Request/Response
type SimilarityRequest struct {
	App1 string `json:"app1"`
	App2 string `json:"app2"`
}

type SimilarityResponse struct {
	SimilarityScore float64 `json:"similarity_score"`
	JaccardIndex    float64 `json:"jaccard_index"`
	LLMScore        float64 `json:"llm_score"`
	Reasoning       string  `json:"reasoning"`
}

// Mock Implementation Methods
func (m *MockAIService) GenerateJDRecommendations(ctx context.Context, req *JDRecommendationRequest) (*JDRecommendationResponse, error) {
	// Simulate AI processing time
	time.Sleep(100 * time.Millisecond)

	// Mock recommendations based on job title and description
	recommendations := []SoftwareRecommendation{
		{
			AppName:        "GitHub Enterprise",
			Category:       "DevOps",
			Tier:           "Pro License",
			RelevanceScore: 95.0,
			Cost:           12000.0,
			Rationale:      "Essential for version control and CI/CD pipelines",
			Features:       []string{"Private repositories", "Advanced security", "Team management"},
			Alternatives:   []string{"GitLab", "Bitbucket"},
		},
		{
			AppName:        "Slack Pro",
			Category:       "Collaboration",
			Tier:           "Pro License",
			RelevanceScore: 88.0,
			Cost:           8500.0,
			Rationale:      "Team communication and file sharing",
			Features:       []string{"Unlimited message history", "File sharing", "Integrations"},
			Alternatives:   []string{"Microsoft Teams", "Discord"},
		},
		{
			AppName:        "Docker Pro",
			Category:       "DevOps",
			Tier:           "Pro License",
			RelevanceScore: 82.0,
			Cost:           7200.0,
			Rationale:      "Containerization and deployment",
			Features:       []string{"Container registry", "Security scanning", "Multi-platform"},
			Alternatives:   []string{"Podman", "LXC"},
		},
	}

	return &JDRecommendationResponse{
		Recommendations: recommendations,
		Confidence:      92.5,
		ProcessingTime:  100 * time.Millisecond,
	}, nil
}

func (m *MockAIService) GenerateConsolidationMemo(ctx context.Context, req *ConsolidationMemoRequest) (*ConsolidationMemoResponse, error) {
	// Simulate AI processing time
	time.Sleep(200 * time.Millisecond)

	memo := fmt.Sprintf(`
# Software Consolidation Analysis

## Executive Summary
Based on analysis of %d applications across %d subsidiaries, we recommend consolidating to achieve %.0f%% cost savings.

## Current State
- Total current cost: ฿%.0f/year
- Number of applications: %d
- Affected subsidiaries: %d

## Recommended Actions
1. Consolidate similar applications
2. Negotiate volume discounts
3. Implement centralized management

## Risk Assessment
Low to medium risk with proper change management.

## Timeline
Implementation can begin immediately with 3-month rollout plan.
`, len(req.Apps), len(req.Subsidiaries), req.TargetSavings/req.CurrentCost*100, req.CurrentCost, len(req.Apps), len(req.Subsidiaries))

	return &ConsolidationMemoResponse{
		Memo:           memo,
		Savings:        req.TargetSavings,
		RiskAssessment: "Low to Medium",
		Timeline:       "3 months",
		Confidence:     87.5,
	}, nil
}

func (m *MockAIService) CalculateSoftwareSimilarity(ctx context.Context, req *SimilarityRequest) (*SimilarityResponse, error) {
	// Simulate AI processing time
	time.Sleep(50 * time.Millisecond)

	// Mock similarity calculation
	similarityScore := 85.0
	jaccardIndex := 0.75
	llmScore := 0.90

	reasoning := fmt.Sprintf("Applications '%s' and '%s' show high functional overlap in core features. Both provide similar capabilities with slight variations in user interface and integration options.", req.App1, req.App2)

	return &SimilarityResponse{
		SimilarityScore: similarityScore,
		JaccardIndex:    jaccardIndex,
		LLMScore:        llmScore,
		Reasoning:       reasoning,
	}, nil
}

func (m *MockAIService) GetAppRecommendationScores(ctx context.Context, prompt string) (string, error) {
	// Mock response - in production, this would call OpenAI API
	// For now, return a mock JSON response
	mockResponse := `[
		{"app_name": "Microsoft Teams", "score": 85},
		{"app_name": "Slack", "score": 80},
		{"app_name": "Zoom", "score": 75},
		{"app_name": "Google Workspace", "score": 70},
		{"app_name": "Jira", "score": 65},
		{"app_name": "Confluence", "score": 60},
		{"app_name": "Figma", "score": 55},
		{"app_name": "Adobe Creative Cloud", "score": 45}
	]`

	return mockResponse, nil
}

func (m *MockAIService) GenerateShortDescription(ctx context.Context, appName, category string, score float64) (string, error) {
	// Mock AI-generated short descriptions based on app type and score
	descriptions := map[string]string{
		"Collaboration":      "Team communication and collaboration platform",
		"Project Management": "Agile project tracking and team coordination",
		"Design":             "Professional design and prototyping tool",
		"Productivity":       "Comprehensive productivity and office suite",
		"Video Conferencing": "Video meetings and online collaboration",
		"Development":        "Software development and version control",
		"Analytics":          "Data analytics and business intelligence",
		"Security":           "Enterprise security and compliance",
		"HR":                 "Human resources management system",
		"Finance":            "Financial management and accounting",
	}

	// Get base description
	desc, ok := descriptions[category]
	if !ok {
		desc = "Professional software application"
	}

	// Enhance based on score
	if score >= 80 {
		desc = "Essential " + desc
	} else if score >= 60 {
		desc = "Recommended " + desc
	}

	return desc, nil
}

// Real AI Service Implementation
type RealAIService struct {
	client *openai.Client
	model  string
}

func NewRealAIService(apiKey string) *RealAIService {
	client := openai.NewClientWithConfig(openai.DefaultConfig(apiKey))

	return &RealAIService{
		client: client,
		model:  openai.GPT4oMini, // or openai.GPT4o, openai.GPT35Turbo
	}
}

// Implementation of real AI methods
func (r *RealAIService) GenerateJDRecommendations(ctx context.Context, req *JDRecommendationRequest) (*JDRecommendationResponse, error) {
	startTime := time.Now()

	prompt := fmt.Sprintf(`You are a software license recommendation expert. Based on the following job description, recommend appropriate software licenses.

Job Title: %s
Department: %s
Company: %s
Experience Level: %s
Skills: %v
Job Description: %s

Provide 3-5 software recommendations with:
- app_name: Name of the application
- category: Software category
- tier: License tier (e.g., "Pro License", "Enterprise")
- relevance_score: Score from 0-100
- cost: Estimated annual cost in THB
- rationale: Brief explanation why this is needed
- features: List of key features needed
- alternatives: List of alternative software

Return the response as a JSON array with this structure:
[
  {
    "app_name": "string",
    "category": "string",
    "tier": "string",
    "relevance_score": number,
    "cost": number,
    "rationale": "string",
    "features": ["string"],
    "alternatives": ["string"]
  }
]`,
		req.JobTitle,
		req.Department,
		req.CompanyCode,
		req.Experience,
		req.Skills,
		req.JobDescription,
	)

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: r.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a software license recommendation expert. Respond only with valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var recommendations []SoftwareRecommendation
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &recommendations); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return &JDRecommendationResponse{
		Recommendations: recommendations,
		Confidence:      85.0, // You could calculate this based on the response
		ProcessingTime:  time.Since(startTime),
	}, nil
}

func (r *RealAIService) GenerateConsolidationMemo(ctx context.Context, req *ConsolidationMemoRequest) (*ConsolidationMemoResponse, error) {
	prompt := fmt.Sprintf(`You are a software consolidation expert. Create an executive memo for consolidating software licenses across subsidiaries.

Current Situation:
- Applications: %v
- Subsidiaries: %v
- Current Annual Cost: ฿%.2f
- Target Savings: ฿%.2f
- Context: %s

Create a professional executive memo that includes:
1. Executive Summary
2. Current State Analysis
3. Recommended Actions
4. Risk Assessment
5. Implementation Timeline
6. Expected Savings

Format the response as a markdown document.`,
		req.Apps,
		req.Subsidiaries,
		req.CurrentCost,
		req.TargetSavings,
		req.Context,
	)

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: r.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a software consolidation expert. Create professional executive memos.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	memo := resp.Choices[0].Message.Content

	return &ConsolidationMemoResponse{
		Memo:           memo,
		Savings:        req.TargetSavings,
		RiskAssessment: "AI-generated assessment included in memo",
		Timeline:       "3-6 months",
		Confidence:     80.0,
	}, nil
}

func (r *RealAIService) CalculateSoftwareSimilarity(ctx context.Context, req *SimilarityRequest) (*SimilarityResponse, error) {
	prompt := fmt.Sprintf(`You are a software analysis expert. Compare the functional similarity between two software applications.

Application 1: %s
Application 2: %s

Provide a detailed similarity analysis including:
- similarity_score: Overall similarity (0-100)
- jaccard_index: Jaccard similarity index (0-1)
- llm_score: Your assessed functional overlap (0-1)
- reasoning: Detailed explanation of similarities and differences

Return the response as JSON:
{
  "similarity_score": number,
  "jaccard_index": number,
  "llm_score": number,
  "reasoning": "string"
}`,
		req.App1,
		req.App2,
	)

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: r.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a software analysis expert. Respond only with valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.5,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var result SimilarityResponse
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return &result, nil
}

func (r *RealAIService) GetAppRecommendationScores(ctx context.Context, prompt string) (string, error) {
	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: r.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a software recommendation expert. Analyze job profiles and recommend appropriate software with relevance scores. Respond only with valid JSON array.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

func (r *RealAIService) GenerateShortDescription(ctx context.Context, appName, category string, score float64) (string, error) {
	prompt := fmt.Sprintf(`Generate a concise, professional description (max 10 words) for this software application:

Application: %s
Category: %s
Relevance Score: %.0f/100

The description should be compelling and highlight the key value proposition. Just return the description text, nothing else.`,
		appName,
		category,
		score,
	)

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: r.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a technical writer creating concise software descriptions. Respond only with the description text.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   50,
		},
	)

	if err != nil {
		// Fallback to category-based description if API fails
		return fmt.Sprintf("%s tool", category), nil
	}

	if len(resp.Choices) == 0 {
		return fmt.Sprintf("%s tool", category), nil
	}

	return resp.Choices[0].Message.Content, nil
}
