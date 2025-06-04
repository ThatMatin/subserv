package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/thatmatin/subserv/pkg/model"
	"github.com/thatmatin/subserv/pkg/repo"
	"gorm.io/gorm"
)

type ProductService interface {
	Get(context.Context, uint) (*model.Product, error)
    GetAll(context.Context) ([]model.Product, error)
}

type productService struct {
    repo repo.ProductRepository
}

func NewProductService(repo repo.ProductRepository) ProductService {
    return &productService { repo: repo }
}

func (s *productService) Get(ctx context.Context, ID uint) (*model.Product, error) {
    product, err := s.repo.GetByID(ctx, ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, err
        } 
        return nil, fmt.Errorf("failed to fetch product: %w", err)
    }

    return product, nil
}

func (s *productService) GetAll(ctx context.Context) ([]model.Product, error) {
    products, err := s.repo.GetAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch all products: %w", err)
    }

    return products, nil
}
