package dto

type SubscriptionUri struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

type CreateSubscriptionRequest struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
}
