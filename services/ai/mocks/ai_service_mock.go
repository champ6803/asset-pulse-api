package mocks

import (
	"asset-pulse-api/services/ai"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockAIService is a mock implementation of ai.AIService
type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) GenerateJDRecommendations(ctx context.Context, req *ai.JDRecommendationRequest) (*ai.JDRecommendationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*ai.JDRecommendationResponse), args.Error(1)
}

func (m *MockAIService) GenerateConsolidationMemo(ctx context.Context, req *ai.ConsolidationMemoRequest) (*ai.ConsolidationMemoResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*ai.ConsolidationMemoResponse), args.Error(1)
}

func (m *MockAIService) CalculateSoftwareSimilarity(ctx context.Context, req *ai.SimilarityRequest) (*ai.SimilarityResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*ai.SimilarityResponse), args.Error(1)
}
