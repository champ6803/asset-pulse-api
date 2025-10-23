package handler

import (
	"asset-pulse-api/services/ai"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"

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
