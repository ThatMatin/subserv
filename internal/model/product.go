package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name     string        `gorm:"not null;size:255"`
	Price    int           `gorm:"not null;type:int"`     // price in cents, e.g., 1999 for $19.99
	TaxRate  uint8         `gorm:"not null;type:tinyint"` // percentage, e.g., 20 for 20%
	Duration time.Duration `gorm:"not null;type:bigint"`  // duration in seconds, e.g., 2592000 for 30 days
}
