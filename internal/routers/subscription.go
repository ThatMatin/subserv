package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/controller"
	"github.com/thatmatin/subserv/internal/middleware"
)

func RegisterSubscriptionRoutes(r *gin.Engine, s *controller.SubscriptionController) {
	subscriptions := r.Group("/subscriptions", middleware.AuthMiddleware())
	{
		subscriptions.GET("/:id", s.GetSubscriptionByID)
		subscriptions.POST("", s.CreateSubscription)
		subscriptions.POST("/:id/purchase", s.Purchase)
		subscriptions.PATCH("/:id/pause", s.PauseSubscription)
		subscriptions.PATCH("/:id/unpause", s.UnpauseSubscription)
		subscriptions.PATCH("/:id/cancel", s.CancelSubscription)
	}
}
