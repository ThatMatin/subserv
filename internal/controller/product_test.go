package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mocklib "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thatmatin/subserv/internal/mock"
	"github.com/thatmatin/subserv/internal/model"
	"github.com/thatmatin/subserv/internal/service"
	"gorm.io/gorm"
)

func TestProductController(t *testing.T) {
	router := gin.Default()

	mockProductRepo := new(mock.MockProductRepo)
	mockProductService := service.NewProductService(mockProductRepo)
	productController := NewProductController(&mockProductService)

	router.GET("/products", productController.GetAllProducts)
	router.GET("/products/:id", productController.GetProductByID)

	t.Run("get all products", func(t *testing.T) {
		mockProductRepo.On("GetAll", mocklib.Anything).Return([]model.Product{
			{Model: gorm.Model{ID: 1}, Name: "Product 1"},
			{Model: gorm.Model{ID: 2}, Name: "Product 2"},
		}, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("get product by id", func(t *testing.T) {
		mockProductRepo.On("GetByID", mocklib.Anything, uint(1)).Return(&model.Product{Model: gorm.Model{ID: 1}}, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		mockProductRepo.AssertExpectations(t)
	})
}
