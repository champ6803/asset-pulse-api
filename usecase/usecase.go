package usecase

import (
	dbRepo "asset-pulse-api/repositories/database"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase/models"
	"context"
)

type useCase struct {
	dbRepo    dbRepo.DatabaseRepository
	aiService ai.AIService
}

//go:generate mockery --name=Usecase
type Usecase interface {
	// User management
	GetUsers(ctx context.Context, in *models.GetUsersInp) (*models.GetUsersResp, error)

	// Authentication
	AuthenticateUser(ctx context.Context, in *models.AuthenticateUserInp) (*models.UserResponse, error)
	HashPassword(password string) (string, error)

	// AI Recommendations (Feature 3)
	GenerateJDRecommendations(ctx context.Context, in *models.JDRecommendationRequest) (*models.JDRecommendationResponse, error)

	// Consolidation (Feature 2)
	GenerateConsolidationMemo(ctx context.Context, in *models.ConsolidationMemoRequest) (*models.ConsolidationMemoResponse, error)

	// Similarity (Feature 1)
	CalculateSoftwareSimilarity(ctx context.Context, in *models.SimilarityRequest) (*models.SimilarityResponse, error)
	GetCrossSubMatch(ctx context.Context, in *models.CrossSubMatchRequest) (*models.CrossSubMatchResponse, error)

	// Seat Optimization (Feature 5)
	GetSeatOptimization(ctx context.Context, in *models.SeatOptimizationRequest) (*models.SeatOptimizationResponse, error)

	// Purchase Templates (Feature 4)
	CreatePurchaseTemplate(ctx context.Context, in *models.PurchaseTemplateRequest) (*models.PurchaseTemplateResponse, error)

	// User Licenses
	GetUserLicenses(ctx context.Context, userID int64, search, status string) (*models.GetUserLicensesResponse, error)
}

type UsecaseOptions struct {
	DBRepo    dbRepo.DatabaseRepository
	AIService ai.AIService
}

func New(options UsecaseOptions) *useCase {
	return &useCase{
		dbRepo:    options.DBRepo,
		aiService: options.AIService,
	}
}
