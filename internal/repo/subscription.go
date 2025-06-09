package repo

import (
	"context"

	"github.com/thatmatin/subserv/internal/model"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	GetByID(ctx context.Context, ID uint) (*model.Subscription, error)
	Create(ctx context.Context, sub *model.Subscription) error
	Save(ctx context.Context, sub *model.Subscription) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) GetByID(ctx context.Context, ID uint) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.db.WithContext(ctx).First(&sub, ID).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *model.Subscription) error {
	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return err
	}
	return nil
}

func (r *subscriptionRepository) Save(ctx context.Context, sub *model.Subscription) error {
	if err := r.db.WithContext(ctx).Save(sub).Error; err != nil {
		return err
	}
	return nil
}
