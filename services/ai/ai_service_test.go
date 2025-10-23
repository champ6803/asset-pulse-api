package ai_test

import (
	"asset-pulse-api/services/ai"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockAIService_GenerateJDRecommendations(t *testing.T) {
	service := ai.NewMockAIService()

	t.Run("Valid Request", func(t *testing.T) {
		req := &ai.JDRecommendationRequest{
			JobTitle:       "Software Engineer",
			JobDescription: "Software Engineer with React and Node.js experience",
			Department:     "Engineering",
			CompanyCode:    "SCB",
			Experience:     "3-5 years",
			Skills:         []string{"React", "Node.js", "TypeScript"},
		}

		resp, err := service.GenerateJDRecommendations(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Recommendations)
		assert.Greater(t, len(resp.Recommendations), 0)

		// Check first recommendation structure
		rec := resp.Recommendations[0]
		assert.NotEmpty(t, rec.AppName)
		assert.NotEmpty(t, rec.Rationale)
		assert.Greater(t, rec.RelevanceScore, 0.0)
		assert.LessOrEqual(t, rec.RelevanceScore, 100.0)
	})

	t.Run("Empty Job Description", func(t *testing.T) {
		req := &ai.JDRecommendationRequest{
			JobTitle:       "",
			JobDescription: "",
			Department:     "Engineering",
			CompanyCode:    "SCB",
		}

		resp, err := service.GenerateJDRecommendations(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Recommendations)
	})
}

func TestMockAIService_GenerateConsolidationMemo(t *testing.T) {
	service := ai.NewMockAIService()

	t.Run("Valid Request", func(t *testing.T) {
		req := &ai.ConsolidationMemoRequest{
			Apps:          []string{"Microsoft Office", "Google Workspace"},
			Subsidiaries:  []string{"SCB", "SCBX"},
			CurrentCost:   10000.0,
			TargetSavings: 2000.0,
			Context:       "Office productivity tools consolidation",
		}

		resp, err := service.GenerateConsolidationMemo(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Memo)
		assert.Greater(t, resp.Savings, 0.0)
		assert.Greater(t, resp.Confidence, 0.0)
	})

	t.Run("Empty Apps List", func(t *testing.T) {
		req := &ai.ConsolidationMemoRequest{
			Apps:          []string{},
			Subsidiaries:  []string{"SCB"},
			CurrentCost:   0.0,
			TargetSavings: 0.0,
		}

		resp, err := service.GenerateConsolidationMemo(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Memo)
	})
}

func TestMockAIService_CalculateSoftwareSimilarity(t *testing.T) {
	service := ai.NewMockAIService()

	t.Run("Valid Request", func(t *testing.T) {
		req := &ai.SimilarityRequest{
			App1: "Microsoft Office",
			App2: "Google Workspace",
		}

		resp, err := service.CalculateSoftwareSimilarity(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.SimilarityScore, 0.0)
		assert.LessOrEqual(t, resp.SimilarityScore, 100.0)
		assert.Greater(t, resp.JaccardIndex, 0.0)
		assert.LessOrEqual(t, resp.JaccardIndex, 1.0)
		assert.NotEmpty(t, resp.Reasoning)
	})

	t.Run("Empty App Names", func(t *testing.T) {
		req := &ai.SimilarityRequest{
			App1: "",
			App2: "",
		}

		resp, err := service.CalculateSoftwareSimilarity(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.SimilarityScore, 0.0)
	})
}
