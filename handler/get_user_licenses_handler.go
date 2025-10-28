package handler

import (
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserLicenses(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get user_id from JWT claims (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		err := fmt.Errorf("unauthorized")
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, "User not authenticated")
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Call usecase
	resp, err := h.useCase.GetUserLicenses(ctx, userID.(int64))
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get user licenses error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Transform to DTO
	output := transformer.SuccessResponse(http.StatusOK, *dto.NewGetUserLicensesDTO(resp))
	c.JSON(http.StatusOK, output)
}
