package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thatmatin/subserv/internal/model"
)

type MockSubscriptionRepo struct {
	mock.Mock
}

func (m *MockSubscriptionRepo) GetByID(ctx context.Context, id uint) (*model.Subscription, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepo) Create(ctx context.Context, subscription *model.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepo) Save(ctx context.Context, subscription *model.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}
