package handler

import (
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUsers(c *gin.Context) {
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
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	var companyCodePtr *string
	if companyCode != "" {
		companyCodePtr = &companyCode
	}

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	resp, err := h.useCase.GetUsers(ctx, &models.GetUsersInp{
		CompanyCode: companyCodePtr,
		Status:      statusPtr,
		Page:        page,
		PageSize:    pageSize,
	})

	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get users error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, *dto.NewGetUsersDTO(resp))
	c.JSON(http.StatusOK, output)
}
