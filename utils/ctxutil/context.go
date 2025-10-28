package ctxutil

import (
	"asset-pulse-api/enums"
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserIDNotFound      = errors.New("user ID not found in context")
	ErrCompanyCodeNotFound = errors.New("company code not found in context")
	ErrRoleNotFound        = errors.New("role not found in context")
)

// GetUserID retrieves the user ID from the Gin context
func GetUserID(c *gin.Context) (int64, error) {
	userID, exists := c.Get(enums.ContextKeyUserID)
	if !exists {
		return 0, ErrUserIDNotFound
	}

	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}

	return id, nil
}

// GetCompanyCode retrieves the company code from the Gin context
func GetCompanyCode(c *gin.Context) (string, error) {
	companyCode, exists := c.Get(enums.ContextKeyCompanyCode)
	if !exists {
		return "", ErrCompanyCodeNotFound
	}

	code, ok := companyCode.(string)
	if !ok {
		return "", errors.New("invalid company code type in context")
	}

	return code, nil
}

// GetRole retrieves the role from the Gin context
func GetRole(c *gin.Context) (string, error) {
	role, exists := c.Get(enums.ContextKeyRole)
	if !exists {
		return "", ErrRoleNotFound
	}

	r, ok := role.(string)
	if !ok {
		return "", errors.New("invalid role type in context")
	}

	return r, nil
}

// GetUsername retrieves the username from the Gin context
func GetUsername(c *gin.Context) (string, error) {
	username, exists := c.Get(enums.ContextKeyUsername)
	if !exists {
		return "", errors.New("username not found in context")
	}

	u, ok := username.(string)
	if !ok {
		return "", errors.New("invalid username type in context")
	}

	return u, nil
}

// GetEmail retrieves the email from the Gin context
func GetEmail(c *gin.Context) (string, error) {
	email, exists := c.Get(enums.ContextKeyEmail)
	if !exists {
		return "", errors.New("email not found in context")
	}

	e, ok := email.(string)
	if !ok {
		return "", errors.New("invalid email type in context")
	}

	return e, nil
}

// GetDepartmentCode retrieves the department code from the Gin context
func GetDepartmentCode(c *gin.Context) (string, error) {
	deptCode, exists := c.Get(enums.ContextKeyDepartmentCode)
	if !exists {
		return "", errors.New("department code not found in context")
	}

	code, ok := deptCode.(string)
	if !ok {
		return "", errors.New("invalid department code type in context")
	}

	return code, nil
}
