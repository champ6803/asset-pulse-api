package dto

import (
	"asset-pulse-api/usecase/models"
	"time"
)

// GetLicensesQueryParams - Query parameters for getting licenses
type GetLicensesQueryParams struct {
	Status      *string `form:"status"`       // active, revoked, expired
	Search      *string `form:"search"`       // search by app name
	Category    *string `form:"category"`     // filter by app category
	LicenseTier *string `form:"license_tier"` // filter by license tier
}

// LicenseDTO - DTO for license item
type LicenseDTO struct {
	LicenseAssignmentID int64      `json:"license_assignment_id"`
	LicenseInventoryID  *int64     `json:"license_inventory_id"`
	AppID               *int64     `json:"app_id"`
	AppName             string     `json:"app_name"`
	AppCategory         *string    `json:"app_category"`
	AppStatus           *string    `json:"app_status"`
	LicenseTier         *string    `json:"license_tier"`
	AssignedAt          time.Time  `json:"assigned_at"`
	TotalSeats          *int       `json:"total_seats"`
	ReservedSeats       *int       `json:"reserved_seats"`
	EffectiveDate       *time.Time `json:"effective_date"`
	ExpireDate          *time.Time `json:"expire_date"`
	IsRevoked           bool       `json:"is_revoked"`
	RevokedAt           *time.Time `json:"revoked_at,omitempty"`
	LastUsed            *time.Time `json:"last_used"`
	UsageFrequency      *int64     `json:"usage_frequency"`
}

// GetLicensesResponse - Response for licenses endpoint
type GetLicensesResponse struct {
	Licenses []LicenseDTO `json:"licenses"`
	Total    int          `json:"total"`
}

// NewGetLicensesDTO - Transform usecase response to DTO
func NewGetLicensesDTO(resp *models.GetLicensesResp) *GetLicensesResponse {
	licenses := make([]LicenseDTO, len(resp.Licenses))
	for i, license := range resp.Licenses {
		licenses[i] = LicenseDTO{
			LicenseAssignmentID: license.LicenseAssignmentID,
			LicenseInventoryID:  license.LicenseInventoryID,
			AppID:               license.AppID,
			AppName:             license.AppName,
			AppCategory:         license.AppCategory,
			AppStatus:           license.AppStatus,
			LicenseTier:         license.LicenseTier,
			AssignedAt:          license.AssignedAt,
			TotalSeats:          license.TotalSeats,
			ReservedSeats:       license.ReservedSeats,
			EffectiveDate:       license.EffectiveDate,
			ExpireDate:          license.ExpireDate,
			IsRevoked:           license.IsRevoked,
			RevokedAt:           license.RevokedAt,
			LastUsed:            license.LastUsed,
			UsageFrequency:      license.UsageFrequency,
		}
	}

	return &GetLicensesResponse{
		Licenses: licenses,
		Total:    resp.Total,
	}
}

// GetActiveLicensesResponse - Response for active licenses endpoint with limit and total
type GetActiveLicensesResponse struct {
	Licenses []LicenseDTO `json:"licenses"`
	Total    int64        `json:"total"`
}

// NewGetLicensesDTO - Transform usecase response to DTO
func NewGetActiveLicensesDTO(resp *models.GetActiveLicensesResp) *GetActiveLicensesResponse {
	licenses := make([]LicenseDTO, len(resp.Licenses))
	for i, license := range resp.Licenses {
		licenses[i] = LicenseDTO{
			LicenseAssignmentID: license.LicenseAssignmentID,
			LicenseInventoryID:  license.LicenseInventoryID,
			AppID:               license.AppID,
			AppName:             license.AppName,
			AppCategory:         license.AppCategory,
			AppStatus:           license.AppStatus,
			LicenseTier:         license.LicenseTier,
			AssignedAt:          license.AssignedAt,
			TotalSeats:          license.TotalSeats,
			ReservedSeats:       license.ReservedSeats,
			EffectiveDate:       license.EffectiveDate,
			ExpireDate:          license.ExpireDate,
			IsRevoked:           license.IsRevoked,
			RevokedAt:           license.RevokedAt,
		}
	}

	return &GetActiveLicensesResponse{
		Licenses: licenses,
		Total:    resp.Total,
	}
}
