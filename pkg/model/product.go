package model

import (
	"time"
)

type Product struct {
	ID       uint
	Name     string
	Price    float32
	TaxRate  uint8
	Duration time.Duration
}
