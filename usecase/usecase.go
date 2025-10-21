package usecase

import (
	dbRepo "asset-pulse-api/repositories/database"
	"asset-pulse-api/usecase/models"
	"asset-pulse-api/utils/apperrs"
	"context"
)

type useCase struct {
	dbRepo dbRepo.DatabaseRepository
}

//go:generate mockery --name=Usecase
type Usecase interface {
	GetUsers(ctx context.Context, in *models.GetUsersInp) (*models.GetUsersResp, *apperrs.AppError)
}

type UsecaseOptions struct {
	DBRepo dbRepo.DatabaseRepository
}

func New(options UsecaseOptions) *useCase {
	return &useCase{
		dbRepo: options.DBRepo,
	}
}

