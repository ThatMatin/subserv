package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/controller"
)

func RegisterProductRoutes(r *gin.Engine, c *controller.ProductController) {
	products := r.Group("/products")
	{
		products.GET("", c.GetAllProducts)
		products.GET("/:id", c.GetProductByID)
	}
}
