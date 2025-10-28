package dto

import (
	"asset-pulse-api/usecase/models"
)

type UserLicenseDTO struct {
	ID               int64    `json:"id"`
	AppID            *int64   `json:"app_id"`
	AppName          *string  `json:"app_name"`
	AppAlias         *string  `json:"app_alias"`
	Category         *string  `json:"category"`
	LicenseTier      *string  `json:"license_tier"`
	AssignedAt       string   `json:"assigned_at"`
	ExpireDate       *string  `json:"expire_date"`
	EffectiveDate    *string  `json:"effective_date"`
	Cost             *float64 `json:"cost"`
	Currency         *string  `json:"currency"`
	UsageCount30Days int      `json:"usage_count_30d"`
	LastUsedAt       *string  `json:"last_used_at"`
}

type GetUserLicensesDTO struct {
	Licenses []UserLicenseDTO `json:"licenses"`
}

func NewGetUserLicensesDTO(resp *models.GetUserLicensesResponse) *GetUserLicensesDTO {
	licenses := make([]UserLicenseDTO, 0, len(resp.Licenses))

	for _, license := range resp.Licenses {
		licenses = append(licenses, UserLicenseDTO{
			ID:               license.ID,
			AppID:            license.AppID,
			AppName:          license.AppName,
			AppAlias:         license.AppAlias,
			Category:         license.Category,
			LicenseTier:      license.LicenseTier,
			AssignedAt:       license.AssignedAt,
			ExpireDate:       license.ExpireDate,
			EffectiveDate:    license.EffectiveDate,
			Cost:             license.Cost,
			Currency:         license.Currency,
			UsageCount30Days: license.UsageCount30Days,
			LastUsedAt:       license.LastUsedAt,
		})
	}

	return &GetUserLicensesDTO{
		Licenses: licenses,
	}
}
