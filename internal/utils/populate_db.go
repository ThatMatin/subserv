package utils

import (
	"github.com/thatmatin/subserv/internal/db"
	"github.com/thatmatin/subserv/internal/model"
)

func PopulateDBWithTestData() {
	db, err := db.Setup()
	if err != nil {
		panic("Failed to setup database: " + err.Error())
	}

	products := []model.Product{
		{Name: "Basic Plan", Price: 999, TaxRate: 15, Duration: 2592000, Description: "Basic plan for individuals"},
		{Name: "Pro Plan", Price: 1999, TaxRate: 5, Duration: 2592000, Description: "Pro plan for small teams"},
		{Name: "Enterprise Plan", Price: 4999, TaxRate: 5, Duration: 2592000, Description: "Enterprise plan with advanced features"},
		{Name: "Premium Plan", Price: 9999, TaxRate: 20, Duration: 2592000, Description: "Premium plan with all features included"},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			panic("Failed to populate database: " + err.Error())
		}
	}

	users := []model.User{
		{Name: "Alice", Email: "alice@d.com"},
		{Name: "Bob", Email: "bob@d.com"},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			panic("Failed to populate database: " + err.Error())
		}
	}
}
