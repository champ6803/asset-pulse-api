package usecase

import (
	"asset-pulse-api/entities"
	dbRepo "asset-pulse-api/repositories/database"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase/models"
	"context"
)

type useCase struct {
	dbRepo              dbRepo.DatabaseRepository
	aiService           ai.AIService
	optimizationService *ai.SeatOptimizationService
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

	// License Management
	GetLicenses(ctx context.Context, in *models.GetLicensesInp) (*models.GetLicensesResp, error)
	GetActiveLicenses(ctx context.Context, in *models.GetActiveLicensesInp) (*models.GetActiveLicensesResp, error)

	// Request Management
	GetPendingRequests(ctx context.Context, in *models.GetPendingRequestsInp) (*models.GetPendingRequestsResp, error)

	// AI App Recommendations
	GetAIRecommendations(ctx context.Context, in models.GetAIRecommendationsInp) (*models.GetAIRecommendationsResp, error)

	// Software license
	GetAllLicenses(ctx context.Context) ([]entities.SoftwareLicense, error)
	InsertLicense(ctx context.Context, license *entities.SoftwareLicense) error
	DeleteLicense(ctx context.Context, id uint) error
	GetGroupedSoftware(ctx context.Context) ([]entities.CurrentGroupedSoftware, error)
	UpsertGroupedSoftware(ctx context.Context, rawJSON []byte) error
	// User Licenses
	GetUserLicenses(ctx context.Context, userID int64, search, status string) (*models.GetUserLicensesResponse, error)
}

type UsecaseOptions struct {
	DBRepo              dbRepo.DatabaseRepository
	AIService           ai.AIService
	OptimizationService *ai.SeatOptimizationService
}

func New(options UsecaseOptions) *useCase {
	return &useCase{
		dbRepo:              options.DBRepo,
		aiService:           options.AIService,
		optimizationService: options.OptimizationService,
	}
}
