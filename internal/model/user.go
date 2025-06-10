package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"not null;type:varchar(100)"`
	Email string `gorm:"not null;unique;type:varchar(100)"`
}
