package model

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID    uint       `gorm:"foreignKey:UserID;type:bigint;not null"`
	ProductID uint       `gorm:"foreignKey:ProductID;type:bigint;not null"`
	State     State      `gorm:"default:0;type:tinyint"` // 0: Pending, 1: Active, 2: Paused, 3: Cancelled, 4: Expired, 5: Failed
	PriceCent int        `gorm:"not null;type:int"`      // price in cents, e.g., 1999 for $19.99
	TaxRate   uint8      `gorm:"default:0;type:tinyint"` // percentage, e.g., 20 for 20%
	Start     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	End       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	PausedAt  *time.Time `gorm:"default:null;type:timestamp"`
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

var StateNames = [...]string{"Pending", "Active", "Paused", "Cancelled", "Expired", "Failed"}
