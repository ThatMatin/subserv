package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/controller"
)

func RegisterSubscriptionRoutes(r *gin.Engine, s *controller.SubscriptionController) {
	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.GET("/:id", s.GetSubscriptionByID)
		subscriptions.POST("", s.CreateSubscription)
		subscriptions.POST("/:id/purchase", s.Purchase)
		subscriptions.PATCH("/:id/pause", s.PauseSubscription)
		subscriptions.PATCH("/:id/unpause", s.UnpauseSubscription)
		subscriptions.PATCH("/:id/cancel", s.CancelSubscription)
	}
}
