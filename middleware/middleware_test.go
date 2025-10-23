package middleware_test

import (
	"asset-pulse-api/middleware"
	"asset-pulse-api/utils/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuth(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Test valid token
	t.Run("Valid Token", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test missing token
	t.Run("Missing Token", func(t *testing.T) {
		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test invalid token
	t.Run("Invalid Token", func(t *testing.T) {
		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRequireRole(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Test valid role
	t.Run("Valid Role", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "manager")

		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.Use(authMiddleware.RequireRole("manager", "cto"))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test invalid role
	t.Run("Invalid Role", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.Use(authMiddleware.RequireRole("manager", "cto"))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRequireCompanyCode(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Test valid company code
	t.Run("Valid Company Code", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.Use(authMiddleware.RequireCompanyCode("SCB", "SCBX"))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test invalid company code
	t.Run("Invalid Company Code", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "OTHER", "Engineering", "employee")

		router := gin.New()
		router.Use(authMiddleware.RequireAuth())
		router.Use(authMiddleware.RequireCompanyCode("SCB", "SCBX"))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
