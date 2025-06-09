package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/dto"
	"github.com/thatmatin/subserv/internal/service"
)

type ProductController struct {
	svc service.ProductService
}

func NewProductController(prodService *service.ProductService) *ProductController {
	controller := &ProductController{
		svc: *prodService,
	}

	return controller
}

// @Summary Get all products
// @Description Fetch all products from the database
// @Tags Products
// @Produce json
// @Success 200 {object} dto.ProductListResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [get]
// @Security BearerAuth
func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.svc.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to fetch products"})
		return
	}
	if len(products) == 0 {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "No products found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductListResponse(products))
}

// @Summary Get product by ID
// @Description Fetch a product by its ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [get]
// @Security BearerAuth
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	var uri dto.ProductRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid product ID"})
		return
	}

	product, err := c.svc.Get(ctx, uri.ID)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusOK, dto.ErrorResponse{Message: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to fetch product"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductResponse(product))
}
