package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/dto"
	"github.com/thatmatin/subserv/internal/service"
)

type SubscriptionController struct {
	svc service.SubscriptionService
}

func NewSubscriptionController(subService *service.SubscriptionService) *SubscriptionController {
	controller := &SubscriptionController{
		svc: *subService,
	}

	return controller
}

// @Summary Get all subscriptions
// @Description Fetch all subscriptions for the authenticated user
// @Tags Subscriptions
// @Accept JSON
// @Produce JSON
// @Success 200 {object} dto.SubscriptionListResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [get]
// @Security BearerAuth
func (c *SubscriptionController) GetSubscriptionByID(ctx *gin.Context) {
	var uri dto.SubscriptionUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := c.svc.Get(ctx, uri.ID)
	if err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription"})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

// @Summary Create a new subscription
// @Description Create a new subscription for the authenticated user
// @Tags Subscriptions
// @Accept JSON
// @Produce JSON
// @Param request body dto.CreateSubscriptionRequest true "Subscription creation request"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
// @Security BearerAuth
func (c *SubscriptionController) CreateSubscription(ctx *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID := userIDVal.(uint)

	subscription, err := c.svc.Create(ctx, req.ProductID, userID)
	if err != nil {
		if err == service.ErrProductNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		if err == service.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err == service.ErrNoPendingPayment {
			ctx.JSON(http.StatusConflict, gin.H{"error": "No pending payment for subscription"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	ctx.JSON(http.StatusCreated, subscription)
}

// @Summary Purchase a subscription
// @Description Purchase a subscription by its ID
// @Tags Subscriptions
// @Accept JSON
// @Produce JSON
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/purchase [post]
// @Security BearerAuth
func (c *SubscriptionController) Purchase(ctx *gin.Context) {
	var uri dto.SubscriptionUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := c.svc.Purchase(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		if err == service.ErrNoPendingPayment {
			ctx.JSON(http.StatusConflict, gin.H{"error": "No pending payment for subscription"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase subscription"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription purchased successfully"})
}

// @Summary Pause a subscription
// @Description Pause a subscription by its ID
// @Tags Subscriptions
// @Accept JSON
// @Produce JSON
// @Param id path string true "Subscription ID"
// @Success 202 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/pause [post]
// @Security BearerAuth
func (c *SubscriptionController) PauseSubscription(ctx *gin.Context) {
	var uri dto.SubscriptionUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := c.svc.Pause(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pause subscription"})
		return
	}

	ctx.Status(http.StatusAccepted)
}

// @Summary Unpause a subscription
// @Description Unpause a subscription by its ID
// @Tags Subscriptions
// @Accept JSON
// @Produce JSON
// @Param id path string true "Subscription ID"
// @Success 202 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/unpause [post]
// @Security BearerAuth
func (c *SubscriptionController) UnpauseSubscription(ctx *gin.Context) {
	var uri dto.SubscriptionUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := c.svc.Unpause(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unpause subscription"})
		return
	}

	ctx.Status(http.StatusAccepted)
}

// @Summary Cancel a subscription
// @Description Cancel a subscription by its ID
// @Tags Subscriptions
// @Accept JSON
// @Produce JSON
// @Param id path string true "Subscription ID"
// @Success 202 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/cancel [post]
// @Security BearerAuth
func (c *SubscriptionController) CancelSubscription(ctx *gin.Context) {
	var uri dto.SubscriptionUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := c.svc.Cancel(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	ctx.Status(http.StatusAccepted)
}
