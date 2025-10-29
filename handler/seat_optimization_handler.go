package handler

import (
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSeatOptimization(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get user context from middleware
	companyCode, _ := c.Get("company_code")
	companyCodeStr := ""
	if companyCode != nil {
		companyCodeStr = companyCode.(string)
	}

	// Get query parameters
	departmentCode := c.Query("department_code")
	appName := c.Query("app_name")
	action := c.Query("action")
	sortBy := c.Query("sort_by")
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")

	// Parse limit and offset
	limitInt := 50
	offsetInt := 0
	if limit != "" {
		fmt.Sscanf(limit, "%d", &limitInt)
	}
	if offset != "" {
		fmt.Sscanf(offset, "%d", &offsetInt)
	}

	// Build request
	req := &models.SeatOptimizationRequest{
		CompanyCode:    companyCodeStr,
		DepartmentCode: departmentCode,
		AppName:        appName,
		Action:         action,
		SortBy:         sortBy,
		Limit:          limitInt,
		Offset:         offsetInt,
	}

	// Call usecase
	resp, err := h.useCase.GetSeatOptimization(ctx, req)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get seat optimization error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, resp)
	c.JSON(http.StatusOK, output)
}
