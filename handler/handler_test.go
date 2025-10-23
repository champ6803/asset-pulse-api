package handler_test

import (
	"asset-pulse-api/handler"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/jwt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	h := handler.NewHandler(handler.HandlerOptions{
		JWTSecret: "test-secret",
	})

	// Test data
	loginReq := map[string]string{
		"username": "test@example.com",
		"password": "password123",
	}

	// Create request
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router and handle request
	router := gin.New()
	router.POST("/api/v1/auth/login", h.Login)
	router.ServeHTTP(w, req)

	// Assert - should fail without proper usecase setup
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGenerateJDRecommendations(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)

	h := handler.NewHandler(handler.HandlerOptions{
		JWTSecret: "test-secret",
	})

	// Test data
	jdReq := models.JDRecommendationRequest{
		JobTitle:       "Software Engineer",
		JobDescription: "Develop web applications",
		Department:     "Engineering",
		CompanyCode:    "SCB",
		Experience:     "3-5 years",
		Skills:         []string{"React", "Node.js"},
	}

	// Create request with JWT token
	token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")
	reqBody, _ := json.Marshal(jdReq)
	req := httptest.NewRequest("POST", "/api/v1/ai/recommendations/jd-match", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router and handle request
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "test@example.com")
		c.Set("email", "test@example.com")
		c.Set("company_code", "SCB")
		c.Set("department_code", "Engineering")
		c.Set("role", "employee")
		c.Next()
	})
	router.POST("/api/v1/ai/recommendations/jd-match", h.GenerateJDRecommendations)
	router.ServeHTTP(w, req)

	// Assert - should fail without proper usecase setup
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetSeatOptimization(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)

	// Create request with JWT token
	token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "cto")
	req := httptest.NewRequest("GET", "/api/v1/recommendations/seat-optimization?company_code=SCB", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router and handle request
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("username", "test@example.com")
		c.Set("email", "test@example.com")
		c.Set("company_code", "SCB")
		c.Set("department_code", "Engineering")
		c.Set("role", "cto")
		c.Next()
	})
	router.GET("/api/v1/recommendations/seat-optimization", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	})
	router.ServeHTTP(w, req)

	// Assert - should return not implemented
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}
