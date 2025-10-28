package usecase

import (
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (u *useCase) AuthenticateUser(ctx context.Context, in *models.AuthenticateUserInp) (*models.UserResponse, error) {
	// Get user by username/email
	user, err := u.dbRepo.GetUserByUsername(ctx, in.Username)
	if err != nil {
		return nil, apperrs.ErrUnauthorized
	}

	// Check if user exists and is active
	if user == nil || (user.Status != nil && *user.Status != "active") {
		return nil, apperrs.ErrUnauthorized
	}

	// Verify password
	if user.Password == nil {
		return nil, apperrs.ErrUnauthorized
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(in.Password))
	if err != nil {
		return nil, apperrs.ErrUnauthorized
	}

	// Get user roles
	roles, err := u.dbRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, apperrs.ErrUnauthorized
	}

	// Determine primary role based on hierarchy
	var primaryRole string
	if len(roles) > 0 {
		// Role hierarchy: Group.CTO > Subsidiary.CTO > Finance.Manager > HR > Employee
		for _, role := range roles {
			switch role {
			case "Group CTO":
				primaryRole = "group-cto"
				break
			case "Subsidiary CTO":
				if primaryRole == "" || primaryRole == "employee" {
					primaryRole = "subsidiary-cto"
				}
			case "Finance Manager":
				if primaryRole == "" || primaryRole == "employee" || primaryRole == "hr" {
					primaryRole = "manager"
				}
			case "HR":
				if primaryRole == "" || primaryRole == "employee" {
					primaryRole = "manager"
				}
			case "Employee":
				if primaryRole == "" {
					primaryRole = "employee"
				}
			}
		}
	}

	// Default to employee if no role found
	if primaryRole == "" {
		primaryRole = "employee"
	}

	// Transform to response
	userResponse := &models.UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		DisplayName:    user.DisplayName,
		CompanyCode:    user.CompanyCode,
		DepartmentCode: user.DepartmentCode,
		Status:         user.Status,
		Role:           &primaryRole,
	}

	return userResponse, nil
}

func (u *useCase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
