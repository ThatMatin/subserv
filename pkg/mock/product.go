package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thatmatin/subserv/pkg/model"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepo) GetAll(ctx context.Context) ([]model.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Product), args.Error(1)
}
