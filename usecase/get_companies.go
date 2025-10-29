package usecase

import (
	"asset-pulse-api/usecase/models"
	"context"
)

func (u *useCase) GetCompanies(ctx context.Context) (*models.GetCompaniesResponse, error) {
	companies, err := u.dbRepo.GetCompanies(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]models.CompanyItem, len(companies))
	for i, c := range companies {
		items[i] = models.CompanyItem{
			ID:   c.ID,
			Code: c.Code,
			Name: c.Name,
		}
	}

	return &models.GetCompaniesResponse{
		Companies: items,
	}, nil
}
