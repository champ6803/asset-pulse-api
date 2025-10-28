package usecase

import (
	"asset-pulse-api/usecase/models"
	"context"
	"fmt"
)

func (u *useCase) GetUserLicenses(ctx context.Context, userID int64, search, status string) (*models.GetUserLicensesResponse, error) {
	// Get license assignments from repository
	assignments, err := u.dbRepo.GetUserLicenseAssignments(ctx, userID, search, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get user license assignments: %w", err)
	}

	// Transform repository results to response models
	licenses := make([]models.UserLicenseResponse, 0, len(assignments))
	for _, assignment := range assignments {
		license := models.UserLicenseResponse{
			ID:               assignment.ID,
			AppID:            assignment.AppID,
			AppName:          assignment.AppName,
			AppAlias:         assignment.AppAlias,
			Category:         assignment.Category,
			LicenseTier:      assignment.LicenseTier,
			AssignedAt:       assignment.AssignedAt,
			ExpireDate:       assignment.ExpireDate,
			EffectiveDate:    assignment.EffectiveDate,
			Cost:             assignment.Cost,
			Currency:         assignment.Currency,
			UsageCount30Days: assignment.UsageCount30Days,
			LastUsedAt:       assignment.LastUsedAt,
		}
		licenses = append(licenses, license)
	}

	return &models.GetUserLicensesResponse{
		Licenses: licenses,
	}, nil
}
