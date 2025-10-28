package models

// UserLicenseResponse represents a user's license with aggregated data
type UserLicenseResponse struct {
	ID               int64    `json:"id"`
	AppID            *int64   `json:"app_id"`
	AppName          *string  `json:"app_name"`
	AppAlias         *string  `json:"app_alias"`
	Category         *string  `json:"category"`
	LicenseTier      *string  `json:"license_tier"`
	AssignedAt       string   `json:"assigned_at"` // RFC3339 format
	ExpireDate       *string  `json:"expire_date"`
	EffectiveDate    *string  `json:"effective_date"`
	Cost             *float64 `json:"cost"`
	Currency         *string  `json:"currency"`
	UsageCount30Days int      `json:"usage_count_30d"`
	LastUsedAt       *string  `json:"last_used_at"` // RFC3339 format
}

// GetUserLicensesResponse represents the response for getting user licenses
type GetUserLicensesResponse struct {
	Licenses []UserLicenseResponse `json:"licenses"`
}
