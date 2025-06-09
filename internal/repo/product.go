package repo

import (
	"context"

	"github.com/thatmatin/subserv/internal/model"
)

type ProductRepository interface {
	GetByID(ctx context.Context, ID uint) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
}
