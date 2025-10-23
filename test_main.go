package main

import (
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase"
	"asset-pulse-api/utils/jwt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "asset_pulse_test")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}

// SetupTestDependencies creates test dependencies
func SetupTestDependencies() (usecase.Usecase, *jwt.JWTManager, ai.AIService) {
	// Create JWT manager
	jwtManager := jwt.NewJWTManager("test-secret-key", 24*time.Hour)

	// Create AI service (mock for testing)
	aiService := ai.NewMockAIService()

	// Create usecase
	uc := usecase.New(usecase.UsecaseOptions{
		AIService: aiService,
	})

	return uc, jwtManager, aiService
}

// CleanupTestData cleans up test data
func CleanupTestData() {
	// Add cleanup logic here if needed
	log.Println("Test cleanup completed")
}
