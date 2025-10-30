package handler

import (
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/ctxutil"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JDRecommendationRequest struct {
	JobTitle       string   `json:"job_title" binding:"required"`
	JobDescription string   `json:"job_description" binding:"required"`
	Department     string   `json:"department" binding:"required"`
	CompanyCode    string   `json:"company_code" binding:"required"`
	Experience     string   `json:"experience"`
	Skills         []string `json:"skills"`
}

type ConsolidationMemoRequest struct {
	Apps          []string `json:"apps" binding:"required"`
	Subsidiaries  []string `json:"subsidiaries" binding:"required"`
	CurrentCost   float64  `json:"current_cost" binding:"required"`
	TargetSavings float64  `json:"target_savings" binding:"required"`
	Context       string   `json:"context"`
}

type SimilarityRequest struct {
	App1 string `json:"app1" binding:"required"`
	App2 string `json:"app2" binding:"required"`
}

type CatalogSearchRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

// NewHireRecommendationsRequest - all optional params per spec
type NewHireRecommendationsRequest struct {
	JobTitle       string   `json:"job_title"`
	JobDescription string   `json:"job_description"`
	Department     string   `json:"department"`
	CompanyCode    string   `json:"company_code"`
	Experience     string   `json:"experience"`
	Skills         []string `json:"skills"`
	AppName        string   `json:"app_name"`
	Limit          int      `json:"limit"`
}

type SimilarAppMeta struct {
	AppName     string  `json:"app_name"`
	Similarity  float64 `json:"similarity"`
	Vendor      string  `json:"vendor"`
	ProductName string  `json:"product_name"`
}

type NewHireRecommendationsResponse struct {
	Recommendations []ai.SoftwareRecommendation `json:"recommendations"`
	Confidence      float64                     `json:"confidence"`
	ProcessingTime  string                      `json:"processing_time"`
	Total           int                         `json:"total_recommendations"`
	Metadata        map[string]interface{}      `json:"metadata"`
}

func (h *Handler) GenerateJDRecommendations(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	var req JDRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		logger.Error(ctx, fmt.Sprintf("Bad request: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Convert to AI service request
	aiReq := &ai.JDRecommendationRequest{
		JobTitle:       req.JobTitle,
		JobDescription: req.JobDescription,
		Department:     req.Department,
		CompanyCode:    req.CompanyCode,
		Experience:     req.Experience,
		Skills:         req.Skills,
	}

	// Call AI service
	response, err := h.aiService.GenerateJDRecommendations(ctx, aiReq)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("AI service error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, response)
	c.JSON(http.StatusOK, output)
}

func (h *Handler) GenerateConsolidationMemo(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	var req ConsolidationMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		logger.Error(ctx, fmt.Sprintf("Bad request: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Convert to AI service request
	aiReq := &ai.ConsolidationMemoRequest{
		Apps:          req.Apps,
		Subsidiaries:  req.Subsidiaries,
		CurrentCost:   req.CurrentCost,
		TargetSavings: req.TargetSavings,
		Context:       req.Context,
	}

	// Call AI service
	response, err := h.aiService.GenerateConsolidationMemo(ctx, aiReq)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("AI service error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, response)
	c.JSON(http.StatusOK, output)
}

func (h *Handler) CalculateSoftwareSimilarity(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	var req SimilarityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		logger.Error(ctx, fmt.Sprintf("Bad request: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Convert to AI service request
	aiReq := &ai.SimilarityRequest{
		App1: req.App1,
		App2: req.App2,
	}

	// Call AI service
	response, err := h.aiService.CalculateSoftwareSimilarity(ctx, aiReq)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("AI service error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, response)
	c.JSON(http.StatusOK, output)
}

// GetAIRecommendations - Get AI-powered app recommendations for user
// @Summary Get AI app recommendations
// @Description Get AI-powered application recommendations based on user's job profile
// @Tags AI Recommendations
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of recommendations (optional)"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/v1/ai/recommendations [get]
func (h *Handler) GetAIRecommendations(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get user_id from context (set by auth middleware)
	userID, err := ctxutil.GetUserID(c)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, fmt.Sprintf("Unauthorized: %v", err))
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Get company_code from context (optional)
	companyCode, _ := ctxutil.GetCompanyCode(c)
	var companyCodePtr *string
	if companyCode != "" {
		companyCodePtr = &companyCode
	}

	// Get limit from query parameter (optional)
	var limitPtr *int
	if limitStr := c.Query("limit"); limitStr != "" {
		var limit int
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && limit > 0 {
			limitPtr = &limit
		}
	}

	// Call usecase
	response, err := h.useCase.GetAIRecommendations(ctx, models.GetAIRecommendationsInp{
		UserID:      userID,
		CompanyCode: companyCodePtr,
		Limit:       limitPtr,
	})
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("AI service error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	output := transformer.SuccessResponse(http.StatusOK, response.Recommendations)
	c.JSON(http.StatusOK, output)
}
func (h *Handler) SearchCatalog(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	var req CatalogSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		logger.Error(ctx, fmt.Sprintf("Bad request: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Convert to catalog search service request
	searchReq := &ai.CatalogSearchRequest{
		Query: req.Query,
		TopK:  req.TopK,
	}

	// Call catalog search service
	response, err := h.catalogSearchService.Search(ctx, searchReq)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Catalog search error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, response)
	c.JSON(http.StatusOK, output)
}

// GetNewHireRecommendations - JD-based recommendations with optional app similarity search
// @Router /api/v1/ai/recommendations/new-hire [post]
func (h *Handler) GetNewHireRecommendations(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// ensure authorized
	if _, err := ctxutil.GetUserID(c); err != nil {
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, fmt.Sprintf("Unauthorized: %v", err))
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	var req NewHireRecommendationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		logger.Error(ctx, fmt.Sprintf("Bad request: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validation: allow either JD fields or app_name to drive the flow
	if strings.TrimSpace(req.AppName) == "" && strings.TrimSpace(req.JobTitle) == "" && strings.TrimSpace(req.JobDescription) == "" {
		res := transformer.ExceptionResponse(http.StatusBadRequest, fmt.Errorf("provide either app_name or both job_title and job_description"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// If app_name provided, switch to catalog search service path immediately
	if strings.TrimSpace(req.AppName) != "" && h.catalogSearchService != nil {
		// Initialize if needed
		if err := h.catalogSearchService.Initialize(ctx); err != nil {
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Catalog init error: %v", err))
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		topK := req.Limit
		if topK <= 0 {
			topK = 5
		}

		searchReq := &ai.CatalogSearchRequest{Query: req.AppName, TopK: topK}
		sResp, sErr := h.catalogSearchService.Search(ctx, searchReq)
		if sErr != nil {
			res := transformer.ExceptionResponse(http.StatusInternalServerError, sErr)
			logger.Error(ctx, fmt.Sprintf("Catalog search error: %v", sErr))
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		// Map catalog search results to SoftwareRecommendation shape
		recs := make([]ai.SoftwareRecommendation, 0, len(sResp.Results))
		for _, r := range sResp.Results {
			// Split functionalities into features by comma
			features := []string{}
			if f := strings.TrimSpace(r.Functionalities); f != "" {
				parts := strings.Split(f, ",")
				for _, p := range parts {
					if s := strings.TrimSpace(p); s != "" {
						features = append(features, s)
					}
				}
			}

			rationale := fmt.Sprintf("Similar to '%s' in category %s. Vendor: %s. Product: %s.", req.AppName, r.Category, r.Vendor, r.ProductName)

			recs = append(recs, ai.SoftwareRecommendation{
				AppName:        r.InputName,
				Category:       r.Category,
				Tier:           "Standard",
				RelevanceScore: r.Similarity, // already scaled 0-100
				Cost:           0,            // unknown from catalog; keep 0
				Rationale:      rationale,
				Features:       features,
				Alternatives:   []string{},
			})
		}

		metadata := map[string]interface{}{
			"job_title":             req.JobTitle,
			"department":            req.Department,
			"company_code":          req.CompanyCode,
			"app_name_used":         true,
			"app_name_searched":     req.AppName,
			"similar_apps_searched": sResp.Results,
		}

		resp := NewHireRecommendationsResponse{
			Recommendations: recs,
			Confidence:      85, // heuristic similar to JD path
			ProcessingTime:  "",
			Total:           len(recs),
			Metadata:        metadata,
		}

		output := transformer.SuccessResponse(http.StatusOK, resp)
		c.JSON(http.StatusOK, output)
		return
	}

	// Fallback to JD-based recommendations when app_name not provided
	aiReq := &ai.JDRecommendationRequest{
		JobTitle:       req.JobTitle,
		JobDescription: req.JobDescription,
		Department:     req.Department,
		CompanyCode:    req.CompanyCode,
		Experience:     req.Experience,
		Skills:         req.Skills,
	}

	aiResp, err := h.aiService.GenerateJDRecommendations(ctx, aiReq)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("AI service error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// If app_name provided, run catalog similarity and enrich results
	metadata := map[string]interface{}{
		"job_title":             req.JobTitle,
		"department":            req.Department,
		"company_code":          req.CompanyCode,
		"app_name_used":         false,
		"similar_apps_searched": nil,
	}

	if strings.TrimSpace(req.AppName) != "" && h.catalogSearchService != nil {
		// Initialize once if needed
		_ = h.catalogSearchService.Initialize(ctx)

		topK := req.Limit
		if topK <= 0 {
			topK = 10
		}
		searchReq := &ai.CatalogSearchRequest{Query: req.AppName, TopK: topK}
		if sResp, sErr := h.catalogSearchService.Search(ctx, searchReq); sErr == nil {
			// Build similarity map by input name
			simMap := map[string]float64{}
			metas := make([]SimilarAppMeta, 0, len(sResp.Results))
			for _, r := range sResp.Results {
				simMap[strings.ToLower(r.InputName)] = r.Similarity
				metas = append(metas, SimilarAppMeta{
					AppName:     r.InputName,
					Similarity:  r.Similarity,
					Vendor:      r.Vendor,
					ProductName: r.ProductName,
				})
			}

			// Attach similarity_score to matching recommendations by name (kept in metadata)
			for i := range aiResp.Recommendations {
				key := strings.ToLower(aiResp.Recommendations[i].AppName)
				if v, ok := simMap[key]; ok {
					_ = v
				}
			}

			metadata["app_name_used"] = true
			metadata["app_name_searched"] = req.AppName
			metadata["similar_apps_searched"] = metas
		} else {
			logger.Error(ctx, fmt.Sprintf("Catalog search error: %v", sErr))
		}
	}

	resp := NewHireRecommendationsResponse{
		Recommendations: aiResp.Recommendations,
		Confidence:      aiResp.Confidence,
		ProcessingTime:  aiResp.ProcessingTime.String(),
		Total:           len(aiResp.Recommendations),
		Metadata:        metadata,
	}

	output := transformer.SuccessResponse(http.StatusOK, resp)
	c.JSON(http.StatusOK, output)
}
