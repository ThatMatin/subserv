package dto

import (
	"time"

	"github.com/thatmatin/subserv/internal/model"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ProductRequest struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

type ProductResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       int           `json:"price"`
	TaxRate     uint8         `json:"tax_rate"`
	Duration    time.Duration `json:"duration" swaggertype:"string,format=duration"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
}

func ToProductResponse(product *model.Product) ProductResponse {
	return ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		TaxRate:  product.TaxRate,
		Duration: product.Duration,
	}
}

func ToProductListResponse(products []model.Product) ProductListResponse {
	res := ProductListResponse{
		Products: make([]ProductResponse, len(products)),
	}

	for i, product := range products {
		res.Products[i] = ToProductResponse(&product)
	}

	return res
}
