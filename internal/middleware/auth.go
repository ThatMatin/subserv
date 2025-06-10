package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/dto"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || token != "Bearer test-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "unauthorized user"})
			return
		}

		// Simulate user ID extraction from token
		c.Set("userID", uint(1))
		c.Next()
	}
}
