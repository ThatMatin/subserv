package repo

import (
	"context"

	"github.com/thatmatin/subserv/internal/model"
)

type SubscriptionRepository interface {
	GetByID(ctx context.Context, ID uint) (*model.Subscription, error)
	Create(ctx context.Context, sub *model.Subscription) error
	Save(ctx context.Context, sub *model.Subscription) error
}
