package dto

type InsertLicenseRequest struct {
	Name              string  `json:"name" binding:"required"`
	ThbPricePerYear   float64 `json:"thbPricePerYear"`
	UsedByCompanyName string  `json:"usedByCompanyName" binding:"required"`
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
	CommonFeatures []string   `json:"commonFeatures"`
}

type Software struct {
	ID                  uint    `json:"id"`
	Name                string  `json:"name"`
	Image               string  `json:"image"`
	LicensePricePerYear float64 `json:"licensePricePerYear"`
}
