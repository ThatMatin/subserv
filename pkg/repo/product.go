package repo

import (
	"context"

	"github.com/thatmatin/subserv/pkg/model"
)

type ProductRepository interface {
	GetByID(ctx context.Context, ID uint) (*model.Product, error)
    GetAll(ctx context.Context) ([]model.Product, error)
}
