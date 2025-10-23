package jwt_test

import (
	"asset-pulse-api/utils/jwt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)

	t.Run("Valid Token Generation", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Empty Secret", func(t *testing.T) {
		jwtManagerEmpty := jwt.NewJWTManager("", 24*time.Hour)
		token, err := jwtManagerEmpty.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		// Empty secret should still work with JWT library, but we can test with invalid secret
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestJWTManager_ValidateToken(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret", 24*time.Hour)

	t.Run("Valid Token", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		claims, err := jwtManager.ValidateToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, int64(1), claims.UserID)
		assert.Equal(t, "test@example.com", claims.Email)
		assert.Equal(t, "test@example.com", claims.Username)
		assert.Equal(t, "SCB", claims.CompanyCode)
		assert.Equal(t, "Engineering", claims.DepartmentCode)
		assert.Equal(t, "employee", claims.Role)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		claims, err := jwtManager.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Empty Token", func(t *testing.T) {
		claims, err := jwtManager.ValidateToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Create a token with very short expiration
		jwtManagerShort := jwt.NewJWTManager("test-secret", 1*time.Millisecond)
		token, _ := jwtManagerShort.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		// Wait for token to expire
		time.Sleep(2 * time.Millisecond)

		claims, err := jwtManagerShort.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestJWTManager_DifferentSecrets(t *testing.T) {
	jwtManager1 := jwt.NewJWTManager("secret1", 24*time.Hour)
	jwtManager2 := jwt.NewJWTManager("secret2", 24*time.Hour)

	t.Run("Token from Different Secret", func(t *testing.T) {
		token, _ := jwtManager1.GenerateToken(1, "test@example.com", "test@example.com", "SCB", "Engineering", "employee")

		claims, err := jwtManager2.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
