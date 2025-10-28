package middleware

import (
	"asset-pulse-api/enums"
	"asset-pulse-api/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *jwt.JWTManager
}

func NewAuthMiddleware(jwtManager *jwt.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(enums.ContextKeyUserID, claims.UserID)
		c.Set(enums.ContextKeyUsername, claims.Username)
		c.Set(enums.ContextKeyEmail, claims.Email)
		c.Set(enums.ContextKeyCompanyCode, claims.CompanyCode)
		c.Set(enums.ContextKeyDepartmentCode, claims.DepartmentCode)
		c.Set(enums.ContextKeyRole, claims.Role)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(enums.ContextKeyRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		c.Abort()
	}
}

func (m *AuthMiddleware) RequireCompanyCode(allowedCompanies ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyCode, exists := c.Get(enums.ContextKeyCompanyCode)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Company code not found",
			})
			c.Abort()
			return
		}

		company := companyCode.(string)
		for _, allowedCompany := range allowedCompanies {
			if company == allowedCompany {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied for this company",
		})
		c.Abort()
	}
}
