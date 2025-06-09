package db

import (
	"fmt"

	"github.com/thatmatin/subserv/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("myproject.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := db.AutoMigrate(&model.Product{}, &model.Subscription{}, &model.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
