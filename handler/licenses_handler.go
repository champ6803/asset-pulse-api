package handler

import (
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/ctxutil"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLicenses(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get user ID from context (set by auth middleware)
	userID, err := ctxutil.GetUserID(c)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, fmt.Sprintf("Failed to get user ID from context: %v", err))
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Get company code from context (optional, for additional filtering)
	companyCode, err := ctxutil.GetCompanyCode(c)
	var companyCodePtr *string
	if err == nil && companyCode != "" {
		companyCodePtr = &companyCode
	}

	// Parse query parameters
	var queryParams dto.GetLicensesQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, fmt.Errorf("invalid query parameters: %v", err))
		logger.Error(ctx, fmt.Sprintf("Invalid query parameters: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Call usecase with authenticated user ID and filters
	resp, err := h.useCase.GetLicenses(ctx, &models.GetLicensesInp{
		UserID:      &userID,
		CompanyCode: companyCodePtr,
		Status:      queryParams.Status,
		Search:      queryParams.Search,
		Category:    queryParams.Category,
		LicenseTier: queryParams.LicenseTier,
	})
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get licenses error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Transform to DTO and return success response
	output := transformer.SuccessResponse(http.StatusOK, *dto.NewGetLicensesDTO(resp))
	c.JSON(http.StatusOK, output)
}

func (h *Handler) GetActiveLicenses(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get user ID from context (set by auth middleware)
	userID, err := ctxutil.GetUserID(c)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, fmt.Sprintf("Failed to get user ID from context: %v", err))
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Get company code from context (optional, for additional filtering)
	companyCode, err := ctxutil.GetCompanyCode(c)
	var companyCodePtr *string
	if err == nil && companyCode != "" {
		companyCodePtr = &companyCode
	}

	// Call usecase with authenticated user ID and limit of 4
	resp, err := h.useCase.GetActiveLicenses(ctx, &models.GetActiveLicensesInp{
		UserID:      &userID,
		CompanyCode: companyCodePtr,
		Limit:       4,
	})
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get active licenses error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Transform to DTO and return success response with total count
	c.JSON(http.StatusOK, transformer.SuccessResponse(http.StatusOK, *dto.NewGetActiveLicensesDTO(resp)))
}
