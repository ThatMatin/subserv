package model

import "time"

type Subscription struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	ProductID uint       `json:"product_id"`
	State     State      `json:"state"`
	Start     time.Time  `json:"start"`
	End       time.Time  `json:"end"`
	PausedAt  *time.Time `json:"paused_at,omitempty"`
	PriceCent int        `json:"price_cent"` // in cents
	TaxRate   uint8      `json:"tax_rate"`   // in percentage
}

type State uint

const (
	Pending State = iota
	Active
	Paused
	Cancelled
	Expired
	Failed
)
