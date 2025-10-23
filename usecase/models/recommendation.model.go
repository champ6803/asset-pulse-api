package models

type JDRecommendationRequest struct {
	JobTitle       string   `json:"job_title" binding:"required"`
	JobDescription string   `json:"job_description" binding:"required"`
	Department     string   `json:"department" binding:"required"`
	CompanyCode    string   `json:"company_code" binding:"required"`
	Experience     string   `json:"experience"`
	Skills         []string `json:"skills"`
}

type JDRecommendationResponse struct {
	Recommendations []SoftwareRecommendation `json:"recommendations"`
	Confidence      float64                  `json:"confidence"`
	ProcessingTime  string                   `json:"processing_time"`
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

type ConsolidationMemoRequest struct {
	Apps          []string `json:"apps" binding:"required"`
	Subsidiaries  []string `json:"subsidiaries" binding:"required"`
	CurrentCost   float64  `json:"current_cost" binding:"required"`
	TargetSavings float64  `json:"target_savings" binding:"required"`
	Context       string   `json:"context"`
}

type ConsolidationMemoResponse struct {
	Memo           string  `json:"memo"`
	Savings        float64 `json:"savings"`
	RiskAssessment string  `json:"risk_assessment"`
	Timeline       string  `json:"timeline"`
	Confidence     float64 `json:"confidence"`
}

type SimilarityRequest struct {
	App1 string `json:"app1" binding:"required"`
	App2 string `json:"app2" binding:"required"`
}

type SimilarityResponse struct {
	SimilarityScore float64 `json:"similarity_score"`
	JaccardIndex    float64 `json:"jaccard_index"`
	LLMScore        float64 `json:"llm_score"`
	Reasoning       string  `json:"reasoning"`
}

type SeatOptimizationRequest struct {
	CompanyCode    string `json:"company_code"`
	DepartmentCode string `json:"department_code"`
	AppName        string `json:"app_name"`
	Action         string `json:"action"` // revoke, reallocate, downgrade
}

type SeatOptimizationResponse struct {
	Optimizations []OptimizationOpportunity `json:"optimizations"`
	TotalSavings  float64                   `json:"total_savings"`
	TotalUsers    int                       `json:"total_users"`
}

type OptimizationOpportunity struct {
	ID               string  `json:"id"`
	AppName          string  `json:"app_name"`
	Department       string  `json:"department"`
	Action           string  `json:"action"`
	InactiveUsers    int     `json:"inactive_users"`
	PendingRequests  int     `json:"pending_requests"`
	CanReallocate    int     `json:"can_reallocate"`
	PotentialSavings float64 `json:"potential_savings"`
	RiskLevel        string  `json:"risk_level"`
	Rationale        string  `json:"rationale"`
}

type PurchaseTemplateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Apps        []string `json:"apps" binding:"required"`
	Cost        float64  `json:"cost"`
	Department  string   `json:"department"`
	CompanyCode string   `json:"company_code"`
}

type PurchaseTemplateResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Apps        []string `json:"apps"`
	Cost        float64  `json:"cost"`
	Department  string   `json:"department"`
	CompanyCode string   `json:"company_code"`
	TimesUsed   int      `json:"times_used"`
	CreatedAt   string   `json:"created_at"`
}

type CrossSubMatchRequest struct {
	CompanyCode string  `json:"company_code"`
	AppName     string  `json:"app_name"`
	Threshold   float64 `json:"threshold"`
}

type CrossSubMatchResponse struct {
	Matches []SimilarApp `json:"matches"`
	Total   int          `json:"total"`
}

type SimilarApp struct {
	AppName         string  `json:"app_name"`
	CompanyCode     string  `json:"company_code"`
	SimilarityScore float64 `json:"similarity_score"`
	JaccardIndex    float64 `json:"jaccard_index"`
	LLMScore        float64 `json:"llm_score"`
	Reasoning       string  `json:"reasoning"`
}
