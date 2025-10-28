package models

import "time"

// GetLicensesInp - Input parameters for getting licenses
type GetLicensesInp struct {
	UserID      *int64  `json:"user_id" form:"user_id"`
	CompanyCode *string `json:"company_code" form:"company_code"`
	Status      *string `json:"status" form:"status"`             // active, revoked, expired
	Search      *string `json:"search" form:"search"`             // search by app name
	Category    *string `json:"category" form:"category"`         // filter by app category
	LicenseTier *string `json:"license_tier" form:"license_tier"` // filter by license tier
}

// LicenseItem - Individual license item in response
type LicenseItem struct {
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

// GetLicensesResp - Response for licenses
type GetLicensesResp struct {
	Licenses []LicenseItem `json:"licenses"`
	Total    int           `json:"total"`
}

// GetActiveLicensesInp - Input parameters for getting active licenses with limit
type GetActiveLicensesInp struct {
	UserID      *int64  `json:"user_id"`
	CompanyCode *string `json:"company_code"`
	Limit       int     `json:"limit"` // default 4
}

// GetActiveLicensesResp - Response for getting active licenses with total count
type GetActiveLicensesResp struct {
	Licenses []LicenseItem `json:"licenses"`
	Total    int64         `json:"total"`
}
