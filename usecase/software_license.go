package usecase

import (
	"asset-pulse-api/entities"
	"context"
	"encoding/json"
	"errors"
)

func (u *useCase) GetAllLicenses(ctx context.Context) ([]entities.SoftwareLicense, error) {
	return u.dbRepo.GetAllSoftwareLicenses(ctx)
}

func (u *useCase) InsertLicense(ctx context.Context, license *entities.SoftwareLicense) error {
	if license.Name == "" || license.UsedByCompanyName == "" {
		return errors.New("name and usedByCompanyName are required")
	}
	return u.dbRepo.InsertSoftwareLicense(ctx, license)
}

func (u *useCase) DeleteLicense(ctx context.Context, id uint) error {
	return u.dbRepo.DeleteSoftwareLicense(ctx, id)
}

func (u *useCase) GetGroupedSoftware(ctx context.Context) ([]entities.CurrentGroupedSoftware, error) {
	return u.dbRepo.GetCurrentGroupedSoftware(ctx)
}

func (u *useCase) UpsertGroupedSoftware(ctx context.Context, rawJSON []byte) error {
	var groups []entities.CurrentGroupedSoftware
	if err := json.Unmarshal(rawJSON, &groups); err != nil {
		return err
	}

	return u.dbRepo.UpsertCurrentGroupedSoftware(ctx, rawJSON)
}
