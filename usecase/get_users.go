package usecase

import (
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"
)

func (u *useCase) GetUsers(ctx context.Context, in *models.GetUsersInp) (*models.GetUsersResp, *apperrs.AppError) {
	// Set default pagination
	page := in.Page
	pageSize := in.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Get users from database
	users, err := u.dbRepo.GetUsers(ctx, in.CompanyCode, in.Status, pageSize, offset)
	if err != nil {
		return nil, apperrs.ErrInternalServer().New().WithCause(err)
	}

	// Get total count
	count, err := u.dbRepo.CountUsers(ctx, in.CompanyCode, in.Status)
	if err != nil {
		return nil, apperrs.ErrInternalServer().New().WithCause(err)
	}

	// Transform to response
	userResponses := make([]models.UserResponse, len(*users))
	for i, user := range *users {
		userResponses[i] = models.UserResponse{
			ID:             user.ID,
			CompanyCode:    user.CompanyCode,
			DepartmentCode: user.DepartmentCode,
			EntraID:        user.EntraID,
			Email:          user.Email,
			DisplayName:    user.DisplayName,
			Title:          user.Title,
			EmployeeID:     user.EmployeeID,
			Status:         user.Status,
		}
	}

	return &models.GetUsersResp{
		Users: userResponses,
		Total: int(count),
	}, nil
}

