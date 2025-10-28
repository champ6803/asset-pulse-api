package dto

type InsertLicenseRequest struct {
	ID                int     `json:"id,omitempty"`
	Name              string  `json:"name" binding:"required"`
	ThbPricePerYear   float64 `json:"thb_price_per_year"`
	UsedByCompanyName string  `json:"used_by_company_name" binding:"required"`
}

type DeleteLicenseRequest struct {
	ID uint `json:"id"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

type GroupedSoftwareResponse struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Items          []Software `json:"items"`
	CommonFeatures []string   `json:"common_features"`
}

type Software struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Image             string  `json:"image"`
	ThbPricePerYear   float64 `json:"thb_price_per_year"`
	UsedByCompanyName string  `json:"used_by_company_name"`
}
