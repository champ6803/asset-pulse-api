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

func (h *Handler) GetPendingRequests(c *gin.Context) {
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
	var query dto.GetPendingRequestsQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, fmt.Errorf("invalid query parameters: %v", err))
		logger.Error(ctx, fmt.Sprintf("Invalid query parameters: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Default limit to 2 if not specified
	if query.Limit == 0 {
		query.Limit = 2
	}

	// Call usecase
	resp, err := h.useCase.GetPendingRequests(ctx, &models.GetPendingRequestsInp{
		UserID:      &userID,
		CompanyCode: companyCodePtr,
		Limit:       query.Limit,
	})
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Failed to get pending requests: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Transform to DTO
	requests := make([]dto.PendingRequestDTO, 0, len(resp.Requests))
	for _, req := range resp.Requests {
		requests = append(requests, dto.PendingRequestDTO{
			ID:        req.ID,
			TicketNo:  req.TicketNo,
			Type:      req.Type,
			Status:    req.Status,
			CreatedAt: req.CreatedAt,
		})
	}

	// Return response
	c.JSON(http.StatusOK, transformer.SuccessResponse(http.StatusOK, dto.GetPendingRequestsResponse{
		Requests: requests,
		Total:    resp.Total,
	}))
}
