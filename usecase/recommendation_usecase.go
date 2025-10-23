package usecase

import (
	"asset-pulse-api/entities"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"
	"time"
)

func (u *useCase) GenerateJDRecommendations(ctx context.Context, in *models.JDRecommendationRequest) (*models.JDRecommendationResponse, error) {
	// Convert to AI service request
	aiReq := &ai.JDRecommendationRequest{
		JobTitle:       in.JobTitle,
		JobDescription: in.JobDescription,
		Department:     in.Department,
		CompanyCode:    in.CompanyCode,
		Experience:     in.Experience,
		Skills:         in.Skills,
	}

	// Call AI service
	response, err := u.aiService.GenerateJDRecommendations(ctx, aiReq)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	// Transform AI response to usecase model
	recommendations := make([]models.SoftwareRecommendation, len(response.Recommendations))
	for i, rec := range response.Recommendations {
		recommendations[i] = models.SoftwareRecommendation{
			AppName:        rec.AppName,
			Category:       rec.Category,
			Tier:           rec.Tier,
			RelevanceScore: rec.RelevanceScore,
			Cost:           rec.Cost,
			Rationale:      rec.Rationale,
			Features:       rec.Features,
			Alternatives:   rec.Alternatives,
		}
	}

	return &models.JDRecommendationResponse{
		Recommendations: recommendations,
		Confidence:      response.Confidence,
		ProcessingTime:  response.ProcessingTime.String(),
	}, nil
}

func (u *useCase) GenerateConsolidationMemo(ctx context.Context, in *models.ConsolidationMemoRequest) (*models.ConsolidationMemoResponse, error) {
	// Convert to AI service request
	aiReq := &ai.ConsolidationMemoRequest{
		Apps:          in.Apps,
		Subsidiaries:  in.Subsidiaries,
		CurrentCost:   in.CurrentCost,
		TargetSavings: in.TargetSavings,
		Context:       in.Context,
	}

	// Call AI service
	response, err := u.aiService.GenerateConsolidationMemo(ctx, aiReq)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	return &models.ConsolidationMemoResponse{
		Memo:           response.Memo,
		Savings:        response.Savings,
		RiskAssessment: response.RiskAssessment,
		Timeline:       response.Timeline,
		Confidence:     response.Confidence,
	}, nil
}

func (u *useCase) CalculateSoftwareSimilarity(ctx context.Context, in *models.SimilarityRequest) (*models.SimilarityResponse, error) {
	// Convert to AI service request
	aiReq := &ai.SimilarityRequest{
		App1: in.App1,
		App2: in.App2,
	}

	// Call AI service
	response, err := u.aiService.CalculateSoftwareSimilarity(ctx, aiReq)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	return &models.SimilarityResponse{
		SimilarityScore: response.SimilarityScore,
		JaccardIndex:    response.JaccardIndex,
		LLMScore:        response.LLMScore,
		Reasoning:       response.Reasoning,
	}, nil
}

func (u *useCase) GetSeatOptimization(ctx context.Context, in *models.SeatOptimizationRequest) (*models.SeatOptimizationResponse, error) {
	// Get optimization opportunities from database
	opportunities, err := u.dbRepo.GetSeatOptimizationOpportunities(ctx, in.CompanyCode, in.DepartmentCode, in.AppName)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	// Transform to response
	optimizations := make([]models.OptimizationOpportunity, len(*opportunities))
	totalSavings := 0.0
	totalUsers := 0

	for i, opp := range *opportunities {
		optimizations[i] = models.OptimizationOpportunity{
			ID:               string(rune(opp.ID)), // Convert int64 to string
			AppName:          opp.AppName,
			Department:       opp.Department,
			Action:           opp.Action,
			InactiveUsers:    opp.InactiveUsers,
			PendingRequests:  opp.PendingRequests,
			CanReallocate:    opp.CanReallocate,
			PotentialSavings: opp.PotentialSavings,
			RiskLevel:        opp.RiskLevel,
			Rationale:        opp.Rationale,
		}
		totalSavings += opp.PotentialSavings
		totalUsers += opp.InactiveUsers + opp.PendingRequests
	}

	return &models.SeatOptimizationResponse{
		Optimizations: optimizations,
		TotalSavings:  totalSavings,
		TotalUsers:    totalUsers,
	}, nil
}

func (u *useCase) CreatePurchaseTemplate(ctx context.Context, in *models.PurchaseTemplateRequest) (*models.PurchaseTemplateResponse, error) {
	// Create template entity
	template := &entities.PurchaseTemplate{
		Name:        in.Name,
		Description: &in.Description,
		CompanyCode: in.CompanyCode,
		Department:  in.Department,
		Apps:        in.Apps[0], // Store first app as string
	}

	// Create template in database
	_, err := u.dbRepo.CreatePurchaseTemplate(ctx, template)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	return &models.PurchaseTemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: *template.Description,
		Apps:        in.Apps,
		Cost:        template.Cost,
		Department:  template.Department,
		CompanyCode: template.CompanyCode,
		TimesUsed:   template.TimesUsed,
		CreatedAt:   template.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *useCase) GetCrossSubMatch(ctx context.Context, in *models.CrossSubMatchRequest) (*models.CrossSubMatchResponse, error) {
	// Get similar apps from database
	similarApps, err := u.dbRepo.GetCrossSubMatch(ctx, in.CompanyCode, in.AppName, in.Threshold)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	// Transform to response
	matches := make([]models.SimilarApp, len(*similarApps))
	for i, app := range *similarApps {
		matches[i] = models.SimilarApp{
			AppName:         app.AppName,
			CompanyCode:     app.CompanyCode,
			SimilarityScore: app.SimilarityScore,
			JaccardIndex:    app.JaccardIndex,
			LLMScore:        app.LLMScore,
			Reasoning:       app.Reasoning,
		}
	}

	return &models.CrossSubMatchResponse{
		Matches: matches,
		Total:   len(matches),
	}, nil
}
