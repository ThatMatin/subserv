package dto

import (
	"time"

	"github.com/thatmatin/subserv/internal/model"
)

type SubscriptionRequest struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

type CreateSubscriptionRequest struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
}

type SubscriptionResponse struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	ProductID uint       `json:"product_id"`
	State     string     `json:"state"`
	PriceCent int        `json:"price_cent"`
	TaxRate   uint8      `json:"tax_rate"`
	Start     time.Time  `json:"start"`
	End       time.Time  `json:"end"`
	PausedAt  *time.Time `json:"paused_at,omitempty"`
}

type SubscriptionMessageResponse struct {
	Message string `json:"message"`
}

func ToSubscriptionResponse(s *model.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:        s.ID,
		ProductID: s.ProductID,
		UserID:    s.UserID,
		State:     model.StateNames[s.State],
		PriceCent: s.PriceCent,
		TaxRate:   s.TaxRate,
		Start:     s.Start,
		End:       s.End,
		PausedAt:  s.PausedAt,
	}
}
