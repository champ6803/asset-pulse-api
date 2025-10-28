package usecase

import (
	"asset-pulse-api/entities"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"
	"fmt"
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
	// Set defaults
	limit := in.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}

	// Call AI service to analyze opportunities
	aiReq := &ai.OptimizationAnalysisRequest{
		CompanyCode:    in.CompanyCode,
		DepartmentCode: in.DepartmentCode,
		AppName:        in.AppName,
		Limit:          limit,
		Offset:         offset,
	}

	aiResponse, err := u.optimizationService.AnalyzeOpportunities(ctx, aiReq)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	// Fetch usage analytics for additional context
	usageAnalytics, err := u.dbRepo.GetLicenseUsageAnalytics(ctx, in.CompanyCode, in.DepartmentCode, in.AppName)
	if err != nil {
		return nil, apperrs.ErrInternal
	}

	// Create a map of app data for lookup
	appDataMap := make(map[string]LicenseUsageData)
	for _, analytic := range usageAnalytics {
		key := fmt.Sprintf("%s-%s", analytic.AppName, analytic.DepartmentCode)
		appDataMap[key] = LicenseUsageData{
			AppID:          analytic.AppID,
			AppName:        analytic.AppName,
			AppCategory:    analytic.AppCategory,
			DepartmentCode: analytic.DepartmentCode,
			DepartmentName: analytic.DepartmentName,
			TotalUsers:     analytic.TotalUsers,
			InactiveUsers:  analytic.InactiveUsers,
		}
	}

	// Transform AI recommendations to optimization opportunities
	optimizations := make([]models.OptimizationOpportunity, 0, len(aiResponse.Recommendations))
	totalSavings := 0.0
	totalUsers := 0
	counts := models.OptimizationCounts{}

	for _, rec := range aiResponse.Recommendations {
		// Get additional data from usage analytics
		key := fmt.Sprintf("%s-%s", rec.AppName, rec.Department)
		appData := appDataMap[key]

		// If not found, try to find any match with the app name
		if appData.AppID == 0 && appData.AppName == "" {
			for _, ad := range usageAnalytics {
				if ad.AppName == rec.AppName {
					appData = LicenseUsageData{
						AppID:          ad.AppID,
						AppName:        ad.AppName,
						AppCategory:    ad.AppCategory,
						DepartmentCode: ad.DepartmentCode,
						DepartmentName: ad.DepartmentName,
						TotalUsers:     ad.TotalUsers,
						InactiveUsers:  ad.InactiveUsers,
					}
					break
				}
			}
		}

		// Fallback if still not found
		if appData.AppCategory == "" {
			appData.AppCategory = "General"
		}
		if appData.AppName == "" {
			appData.AppName = rec.AppName
		}

		opp := models.OptimizationOpportunity{
			ID:                   fmt.Sprintf("opt-%d", len(optimizations)+1),
			AppID:                appData.AppID,
			AppName:              rec.AppName,
			AppCategory:          appData.AppCategory,
			Department:           rec.Department,
			DepartmentCode:       appData.DepartmentCode,
			CompanyCode:          in.CompanyCode,
			Action:               rec.Action,
			InactiveUsers:        appData.InactiveUsers,
			PendingRequests:      0,
			CanReallocate:        0,
			PotentialSavings:     rec.PotentialSavings,
			RiskLevel:            rec.RiskLevel,
			Priority:             rec.Priority,
			Rationale:            rec.Rationale,
			AIGeneratedRationale: rec.Rationale,
			FromDepartment:       rec.FromDepartment,
			ToDepartment:         rec.ToDepartment,
			DowngradeFrom:        rec.DowngradeFrom,
			DowngradeTo:          rec.DowngradeTo,
			LastUsedDays:         90,
		}

		// Update counts
		counts.Total++
		switch rec.Action {
		case "revoke":
			counts.Revoke++
		case "reallocate":
			counts.Reallocate++
		case "downgrade":
			counts.Downgrade++
		}

		totalSavings += rec.PotentialSavings
		totalUsers += appData.InactiveUsers

		// Apply filters
		if in.Action != "" && rec.Action != in.Action {
			continue
		}

		optimizations = append(optimizations, opp)
	}

	return &models.SeatOptimizationResponse{
		Optimizations: optimizations,
		TotalSavings:  totalSavings,
		TotalUsers:    totalUsers,
		Counts:        counts,
	}, nil
}

// LicenseUsageData helper struct
type LicenseUsageData struct {
	AppID          int64
	AppName        string
	AppCategory    string
	DepartmentCode string
	DepartmentName string
	TotalUsers     int
	InactiveUsers  int
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
