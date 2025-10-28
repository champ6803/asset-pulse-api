package entities

import "time"

type SoftwareLicense struct {
	ID                uint    `gorm:"primaryKey"`
	Name              string  `gorm:"size:255;not null"`
	THBPricePerYear   float64 `gorm:"not null"`
	UsedByCompanyName string  `gorm:"size:255;not null"`
	CreatedAt         time.Time
}
