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

// @Summary Get subscription
// @Description Fetch subscription by ID for the authenticated user
// @Tags Subscriptions
// @Produce json
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [get]
// @Security BearerAuth
func (c *SubscriptionController) GetSubscriptionByID(ctx *gin.Context) {
	var uri dto.SubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid subscription ID"})
		return
	}

	subscription, err := c.svc.Get(ctx, uri.ID)
	if err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to fetch subscription"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToSubscriptionResponse(subscription))
}

// @Summary Create a new subscription
// @Description Create a new subscription for the authenticated user
// @Tags Subscriptions
// @Accept json
// @Produce json
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
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request body"})
		return
	}

	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Unauthorized"})
		return
	}

	userID := userIDVal.(uint)

	subscription, err := c.svc.Create(ctx, req.ProductID, userID)
	if err != nil {
		if err == service.ErrProductNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Product not found"})
			return
		}
		if err == service.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "User not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to create subscription"})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToSubscriptionResponse(subscription))
}

// @Summary Purchase a subscription
// @Description Purchase a subscription by its ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/purchase [patch]
// @Security BearerAuth
func (c *SubscriptionController) Purchase(ctx *gin.Context) {
	var uri dto.SubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid subscription ID"})
		return
	}

	if err := c.svc.Purchase(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Subscription not found"})
			return
		}
		if err == service.ErrNoPendingPayment {
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{Message: "No pending payment for this subscription"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to purchase subscription"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SubscriptionMessageResponse{Message: "Subscription purchased successfully"})
}

// @Summary Pause a subscription
// @Description Pause a subscription by its ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 202 {object} dto.SubscriptionMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/pause [patch]
// @Security BearerAuth
func (c *SubscriptionController) PauseSubscription(ctx *gin.Context) {
	var uri dto.SubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid subscription ID"})
		return
	}

	if err := c.svc.Pause(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to pause subscription"})
		return
	}

	ctx.JSON(http.StatusAccepted, dto.SubscriptionMessageResponse{Message: "Subscription paused successfully"})
}

// @Summary Unpause a subscription
// @Description Unpause a subscription by its ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 202 {object} dto.SubscriptionMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/unpause [patch]
// @Security BearerAuth
func (c *SubscriptionController) UnpauseSubscription(ctx *gin.Context) {
	var uri dto.SubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid subscription ID"})
		return
	}

	if err := c.svc.Unpause(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to unpause subscription"})
		return
	}

	ctx.JSON(http.StatusAccepted, dto.SubscriptionMessageResponse{Message: "Subscription unpaused successfully"})
}

// @Summary Cancel a subscription
// @Description Cancel a subscription by its ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 202 {object} dto.SubscriptionMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id}/cancel [patch]
// @Security BearerAuth
func (c *SubscriptionController) CancelSubscription(ctx *gin.Context) {
	var uri dto.SubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid subscription ID"})
		return
	}

	if err := c.svc.Cancel(ctx, uri.ID); err != nil {
		if err == service.ErrSubscriptionNotFound {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Subscription not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to cancel subscription"})
		return
	}

	ctx.JSON(http.StatusAccepted, dto.SubscriptionMessageResponse{Message: "Subscription cancelled successfully"})
}
