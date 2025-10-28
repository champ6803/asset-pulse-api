package usecase

import (
	"asset-pulse-api/entities"
	repositories "asset-pulse-api/repositories/database"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"
	"encoding/json"
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

// AI App Recommendations

type OpenAIAppScore struct {
	AppName string  `json:"app_name"`
	Score   float64 `json:"score"`
}

func (u *useCase) GetAIRecommendations(ctx context.Context, inp models.GetAIRecommendationsInp) (*models.GetAIRecommendationsResp, error) {
	// 1. Get user's job profile
	jobProfile, err := u.dbRepo.GetUserJobProfile(ctx, inp.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user job profile: %w", err)
	}

	// 2. Get all active apps
	apps, err := u.dbRepo.GetActiveApps(ctx, inp.CompanyCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get active apps: %w", err)
	}

	fmt.Println("apps", len(apps))

	if len(apps) == 0 {
		return &models.GetAIRecommendationsResp{
			Recommendations: []models.AIRecommendationItem{},
		}, nil
	}

	// 3. Build prompt for OpenAI
	prompt := buildRecommendationPrompt(jobProfile, apps)

	// 4. Call OpenAI to get scores
	aiScores, err := u.aiService.GetAppRecommendationScores(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI scores: %w", err)
	}

	// 5. Parse AI response and match with apps
	var openAIScores []OpenAIAppScore
	if err := json.Unmarshal([]byte(aiScores), &openAIScores); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// 6. Map scores to apps and add mock cost with AI-generated descriptions
	recommendations, err := mapScoresToAppsWithDesc(ctx, u.dbRepo, apps, openAIScores, inp.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to map scores to apps: %w", err)
	}

	return &models.GetAIRecommendationsResp{
		Recommendations: recommendations,
	}, nil
}

func buildRecommendationPrompt(jobProfile *entities.JobProfile, apps []entities.App) string {
	appList := ""
	for i, app := range apps {
		category := "General"
		if app.Category != nil {
			category = *app.Category
		}
		appList += fmt.Sprintf("%d. %s (Category: %s)\n", i+1, app.Name, category)
	}

	prompt := fmt.Sprintf(`You are an expert IT asset manager. Given a job profile and a list of available applications, score each application based on how relevant it is for this job profile.

Job Profile:
- Name: %s
- Description: %s

Available Applications:
%s

Please provide a JSON array with scores (0-100) for each application, where:
- 0-20: Not relevant
- 21-40: Slightly relevant
- 41-60: Moderately relevant
- 61-80: Very relevant
- 81-100: Essential

Return ONLY a JSON array in this exact format:
[
  {"app_name": "Application Name", "score": 75},
  {"app_name": "Another App", "score": 50}
]

Include all applications in your response.`,
		jobProfile.Name,
		getStringValue(jobProfile.Description),
		appList,
	)

	return prompt
}

func mapScoresToAppsWithDesc(ctx context.Context, dbRepo repositories.DatabaseRepository, apps []entities.App, openAIScores []OpenAIAppScore, limit *int) ([]models.AIRecommendationItem, error) {
	scoreMap := make(map[string]float64)
	for _, score := range openAIScores {
		scoreMap[score.AppName] = score.Score
	}

	var recommendations []models.AIRecommendationItem
	for _, app := range apps {
		score, found := scoreMap[app.Name]
		if !found {
			score = 0 // Default score if not found
		}

		// Get category string
		category := "General"
		if app.Category != nil {
			category = *app.Category
		}

		recommendations = append(recommendations, models.AIRecommendationItem{
			AppID:          app.ID,
			AppName:        app.Name,
			AppLogoURL:     GetMockLogoURL(app.Name),
			AppCategory:    category,
			LicenseTier:    *app.ApplicationTier, // Default tier
			RelevanceScore: score,
			CostPerSeat:    0,  // Will be set later
			Rationale:      "", // Will be generated later
		})
	}

	// Sort by relevance score (highest first)
	for i := 0; i < len(recommendations); i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[j].RelevanceScore > recommendations[i].RelevanceScore {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	// Apply limit if specified
	if limit != nil && *limit > 0 && *limit < len(recommendations) {
		recommendations = recommendations[:*limit]
	}

	// Get app IDs for price lookup
	appIDs := make([]int64, len(recommendations))
	for i, rec := range recommendations {
		appIDs[i] = rec.AppID
	}

	// Fetch prices from database
	priceMap, err := dbRepo.GetAppPrices(ctx, appIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get app prices: %w", err)
	}

	// Now generate cost and rationale only for the limited results
	for i := range recommendations {
		// Get price from database or use default
		price, found := priceMap[recommendations[i].AppID]
		if !found || price == 0 {
			// Default price if not found (50-500 THB)
			price = 0
		}
		recommendations[i].CostPerSeat = price

		// Generate rationale based on score
		recommendations[i].Rationale = generateRationale(
			recommendations[i].RelevanceScore,
			recommendations[i].AppName,
			recommendations[i].AppCategory,
		)
	}

	return recommendations, nil
}

func generateRationale(score float64, appName, category string) string {
	switch {
	case score >= 80:
		return fmt.Sprintf("%s is highly recommended for this role. Essential for %s tasks with excellent feature alignment.", appName, category)
	case score >= 60:
		return fmt.Sprintf("%s is strongly recommended. Very relevant for %s workflows in this position.", appName, category)
	case score >= 40:
		return fmt.Sprintf("%s may be useful for certain %s tasks in this role.", appName, category)
	case score >= 20:
		return fmt.Sprintf("%s has limited relevance but could support occasional %s needs.", appName, category)
	default:
		return fmt.Sprintf("%s has minimal relevance for this job profile.", appName)
	}
}

func getStringValue(s *string) string {
	if s == nil {
		return "N/A"
	}
	return *s
}
