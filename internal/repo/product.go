package repo

import (
	"context"

	"github.com/thatmatin/subserv/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByID(ctx context.Context, ID uint) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetByID(ctx context.Context, ID uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, ID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
