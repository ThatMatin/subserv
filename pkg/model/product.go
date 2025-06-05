package model

import (
	"time"
)

type Product struct {
	ID       uint
	Name     string
	Price    int
	TaxRate  uint8
	Duration time.Duration
}
