package handler

import (
	"asset-pulse-api/enums"
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
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
			apperr := apperrs.ErrInternalServer().New().WithCause(fmt.Errorf("%v", r))
			res := transformer.ExceptionResponse(apperr.Status, apperr)
			logger.Error(ctx, fmt.Sprintf(enums.ErrorMsgTemplate, apperr.Status, apperr.Message, apperr.Error()))
			c.JSON(apperr.Status, res)
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

	resp, appErr := h.useCase.GetUsers(ctx, &models.GetUsersInp{
		CompanyCode: companyCodePtr,
		Status:      statusPtr,
		Page:        page,
		PageSize:    pageSize,
	})

	if appErr != nil {
		res := transformer.ExceptionResponse(appErr.Status, appErr)
		logger.Error(ctx, fmt.Sprintf(enums.ErrorMsgTemplate, appErr.Status, appErr.Message, appErr.Error()))
		c.JSON(appErr.Status, res)
		return
	}

	output := transformer.SuccessResponse(http.StatusOK, *dto.NewGetUsersDTO(resp))
	c.JSON(http.StatusOK, output)
}

