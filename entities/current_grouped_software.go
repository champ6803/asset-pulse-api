package entities

import (
	"time"
)

type RawCurrentGroupedSoftware struct {
	ID        uint   `gorm:"primaryKey"`
	JSONData  []byte `gorm:"type:TEXT;not null"`
	UpdatedAt time.Time
}

type SoftwareItem struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Image             string  `json:"image"`
	ThbPricePerYear   float64 `json:"thb_price_per_year"`
	UsedByCompanyName string  `json:"used_by_company_name"`
}

type CurrentGroupedSoftware struct {
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Items          []SoftwareItem `json:"items"`
	CommonFeatures []string       `json:"common_features"`
}
