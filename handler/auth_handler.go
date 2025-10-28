package handler

import (
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      dto.UserDTO `json:"user"`
	ExpiresAt string      `json:"expires_at"`
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		logger.Error(ctx, fmt.Sprintf("Bad request: %v", err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Authenticate user
	user, err := h.useCase.AuthenticateUser(ctx, &models.AuthenticateUserInp{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, fmt.Sprintf("Authentication error: %v", err))
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Get role from user response
	userRole := "employee" // Default role
	if user.Role != nil {
		userRole = *user.Role
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(
		user.ID,
		*user.Username,
		*user.Email,
		*user.CompanyCode,
		*user.DepartmentCode,
		userRole,
	)

	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Token generation error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	response := LoginResponse{
		Token:     token,
		User:      *dto.NewUserDTO(user),
		ExpiresAt: "24h", // Token expires in 24 hours
	}

	output := transformer.SuccessResponse(http.StatusOK, response)
	c.JSON(http.StatusOK, output)
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	// Get user info from JWT claims
	userID, exists := c.Get("user_id")
	if !exists {
		err := fmt.Errorf("unauthorized")
		res := transformer.ExceptionResponse(http.StatusUnauthorized, err)
		logger.Error(ctx, "User not authenticated")
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	username, _ := c.Get("username")
	email, _ := c.Get("email")
	companyCode, _ := c.Get("company_code")
	departmentCode, _ := c.Get("department_code")
	role, _ := c.Get("role")

	user := &models.UserResponse{
		ID:             userID.(int64),
		Username:       stringPtr(username.(string)),
		Email:          stringPtr(email.(string)),
		CompanyCode:    stringPtr(companyCode.(string)),
		DepartmentCode: stringPtr(departmentCode.(string)),
		Status:         stringPtr("active"),
		Role:           stringPtr(role.(string)),
	}

	output := transformer.SuccessResponse(http.StatusOK, *dto.NewUserDTO(user))
	c.JSON(http.StatusOK, output)
}

func (h *Handler) Logout(c *gin.Context) {
	// For JWT tokens, logout is handled client-side by removing the token
	// In a more sophisticated setup, you might maintain a blacklist of tokens
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
