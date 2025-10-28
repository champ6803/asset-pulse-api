package usecase

import (
	"asset-pulse-api/usecase/models"
	"context"
)

func (u *useCase) GetActiveLicenses(ctx context.Context, in *models.GetActiveLicensesInp) (*models.GetActiveLicensesResp, error) {
	// Set default limit to 4 if not provided
	limit := in.Limit
	if limit == 0 {
		limit = 4
	}

	// Set active status filter
	activeStatus := "active"

	// Get licenses from repository with limit and total count
	results, total, err := u.dbRepo.GetLicensesWithLimit(
		ctx,
		in.UserID,
		in.CompanyCode,
		&activeStatus,
		nil, // search
		nil, // category
		nil, // licenseTier
		limit,
	)
	if err != nil {
		return nil, err
	}

	// Transform entity results to response model (type-safe)
	licenses := make([]models.LicenseItem, 0, len(results))
	for _, result := range results {
		license := models.LicenseItem{
			LicenseAssignmentID: result.LicenseAssignmentID,
			LicenseInventoryID:  result.LicenseInventoryID,
			AppID:               &result.AppID,
			AppName:             result.AppName,
			AppCategory:         result.AppCategory,
			AppStatus:           result.AppStatus,
			AppLogoURL:          GetMockLogoURL(result.AppName),
			LicenseTier:         result.LicenseTier,
			AssignedAt:          result.AssignedAt,
			RevokedAt:           result.RevokedAt,
			IsRevoked:           result.IsRevoked,
			TotalSeats:          result.TotalSeats,
			ReservedSeats:       result.ReservedSeats,
			EffectiveDate:       result.EffectiveDate,
			ExpireDate:          result.ExpireDate,
			LastUsed:            result.LastUsed,
			UsageFrequency:      result.UsageFrequency,
		}

		licenses = append(licenses, license)
	}

	return &models.GetActiveLicensesResp{
		Licenses: licenses,
		Total:    total,
	}, nil
}
