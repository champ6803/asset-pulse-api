package handler

import (
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetConsolidationOpportunities(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get query parameters
	companyCode := c.Query("company_code")

	opportunities, err := h.dbRepo.GetConsolidationOpportunities(ctx, companyCode)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get consolidation opportunities error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, opportunities)
	c.JSON(http.StatusOK, output)
}

func (h *Handler) GetConsolidationOpportunityByID(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, fmt.Errorf("invalid id parameter"))
		logger.Error(ctx, fmt.Sprintf("Invalid id parameter: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	opportunity, err := h.dbRepo.GetConsolidationOpportunityByID(ctx, id)
	if err != nil {
		// Check if record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res := transformer.ExceptionResponse(http.StatusNotFound, fmt.Errorf("consolidation opportunity not found"))
			logger.Error(ctx, fmt.Sprintf("Consolidation opportunity not found for id %d: %v", id, err))
			c.JSON(http.StatusNotFound, res)
			return
		}

		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get consolidation opportunity error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, opportunity)
	c.JSON(http.StatusOK, output)
}
