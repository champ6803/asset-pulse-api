package models

// GetAIRecommendationsInp - Input for AI recommendations
type GetAIRecommendationsInp struct {
	UserID      int64
	CompanyCode *string
	Limit       *int // Optional limit for number of recommendations
}

// AIRecommendationItem - Individual recommendation result
type AIRecommendationItem struct {
	AppID          int64   `json:"app_id"`
	AppName        string  `json:"app_name"`
	AppCategory    string  `json:"app_category"`
	LicenseTier    string  `json:"license_tier"`
	RelevanceScore float64 `json:"relevance_score"`
	CostPerSeat    float64 `json:"cost_per_seat"`
	Rationale      string  `json:"rationale"`
}

// GetAIRecommendationsResp - Response for AI recommendations
type GetAIRecommendationsResp struct {
	Recommendations []AIRecommendationItem `json:"recommendations"`
}
