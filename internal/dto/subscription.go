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
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	State     uint8  `json:"state"`
	PriceCent int    `json:"price_cent"`
	TaxRate   uint8  `json:"tax_rate"`
	Start     string `json:"start"`
	End       string `json:"end"`
	PausedAt  string `json:"paused_at,omitempty"`
}

type SubscriptionMessageResponse struct {
	Message string `json:"message"`
}

func ToSubscriptionResponse(s *model.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:        s.ID,
		ProductID: s.ProductID,
		UserID:    s.UserID,
		State:     uint8(s.State),
		PriceCent: s.PriceCent,
		TaxRate:   s.TaxRate,
		Start:     s.Start.Format(time.RFC3339),
		End:       s.End.Format(time.RFC3339),
		PausedAt:  formatPausedAt(s.PausedAt),
	}
}

func formatPausedAt(pausedAt *time.Time) string {
	if pausedAt == nil {
		return ""
	}
	return pausedAt.Format(time.RFC3339)
}
