package model

import "time"

type Subscription struct {
	ID        uint
	UserID    uint
	ProductID uint
	Start     time.Time
	End       time.Time
	PausedAt  *time.Time
	State     State
	PriceCent int
	TaxRate   uint8
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
