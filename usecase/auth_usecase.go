package usecase

import (
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"
	"fmt"

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
	fmt.Println("Password_db: ", *user.Password)
	fmt.Println("Password_input: ", in.Password)
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(in.Password))
	if err != nil {
		return nil, apperrs.ErrUnauthorized
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
