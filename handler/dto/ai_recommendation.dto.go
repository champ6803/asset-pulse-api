package dto

// AIRecommendationDTO - Response DTO for individual recommendation
type AIRecommendationDTO struct {
	AppID          int64   `json:"app_id"`
	AppName        string  `json:"app_name"`
	AppCategory    string  `json:"app_category"`
	LicenseTier    string  `json:"license_tier"`
	RelevanceScore float64 `json:"relevance_score"`
	CostPerSeat    float64 `json:"cost_per_seat"`
	Rationale      string  `json:"rationale"`
}
