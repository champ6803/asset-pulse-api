package usecase

import (
	"asset-pulse-api/usecase/models"
	"context"
)

func (u *useCase) GetLicenses(ctx context.Context, in *models.GetLicensesInp) (*models.GetLicensesResp, error) {

	// Get licenses from repository with filters
	results, err := u.dbRepo.GetLicenses(
		ctx,
		in.UserID,
		in.CompanyCode,
		in.Status,
		in.Search,
		in.Category,
		in.LicenseTier,
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

	return &models.GetLicensesResp{
		Licenses: licenses,
	}, nil
}
