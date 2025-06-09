package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/controller"
)

func registerProductRoutes(r *gin.Engine, c *controller.ProductController) *gin.Engine {
	products := r.Group("/products")
	{
		products.GET("", c.GetProducts)
		products.GET("/:id", c.GetProductByID)
	}

	return r
}
