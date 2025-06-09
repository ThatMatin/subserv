package model

import (
	"time"
)

type Product struct {
	ID       uint          `json:"id"`
	Name     string        `json:"name"`
	Price    int           `json:"price"`    // in cents
	TaxRate  uint8         `json:"tax_rate"` // in percentage
	Duration time.Duration `json:"duration"`
}
