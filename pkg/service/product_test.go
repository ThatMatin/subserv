package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thatmatin/subserv/pkg/mock"
	"github.com/thatmatin/subserv/pkg/model"
	"gorm.io/gorm"
)

func TestFetchSingleProduct(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name             string
		inputID          uint
		expectedName     string
		expectedErr      error
		errorContains    string
		expectedDuration time.Duration
		setupMock        func(repo *mock.MockProductRepo)
	}{
		{
			name:             "product exists",
			inputID:          1,
			expectedName:     "flowmotion",
			expectedErr:      nil,
			expectedDuration: time.Hour * 24 * 30,
			setupMock: func(repo *mock.MockProductRepo) {
				product := &model.Product{ID: uint(1), Name: "flowmotion", Duration: time.Hour * 24 * 30}
				repo.On("GetByID", ctx, uint(1)).Return(product, nil)
			},
		},
		{
			name:        "product not found",
			inputID:     2,
			expectedErr: gorm.ErrRecordNotFound,
			setupMock: func(repo *mock.MockProductRepo) {
				repo.On("GetByID", ctx, uint(2)).Return((*model.Product)(nil), gorm.ErrRecordNotFound)
			},
		},
		{
			name:          "any other repository error",
			inputID:       3,
			expectedErr:   errors.New(""),
			errorContains: "failed to fetch",
			setupMock: func(repo *mock.MockProductRepo) {
				repo.On("GetByID", ctx, uint(3)).Return((*model.Product)(nil), gorm.ErrInvalidDB)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProductRepo := new(mock.MockProductRepo)
			tc.setupMock(mockProductRepo)

			svc := NewProductService(mockProductRepo)

			product, err := svc.Get(ctx, tc.inputID)
			if tc.expectedErr != nil {
				require.Error(t, err)
				if tc.errorContains == "" {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.ErrorContains(t, err, tc.errorContains)
				}

			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedName, product.Name)
				require.Equal(t, tc.expectedDuration, product.Duration)
			}

			mockProductRepo.AssertExpectations(t)
		})
	}
}

func TestFetchListProduct(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name          string
		expectedLen   int
		errorContains string
		setupMock     func(repo *mock.MockProductRepo)
	}{
		{
			name:        "returns product list",
			expectedLen: 2,
			setupMock: func(repo *mock.MockProductRepo) {
				repo.On("GetAll", ctx).Return([]model.Product{
					{ID: 1, Name: "flowmotion"},
					{ID: 2, Name: "flexifit"},
				}, nil)
			},
		},
		{
			name:        "empty list",
			expectedLen: 0,
			setupMock: func(repo *mock.MockProductRepo) {
				repo.On("GetAll", ctx).Return([]model.Product{}, nil)
			},
		},
		{
			name:          "repository error",
			expectedLen:   0,
			errorContains: "failed to fetch",
			setupMock: func(repo *mock.MockProductRepo) {
				repo.On("GetAll", ctx).Return([]model.Product{}, errors.New("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProductRepo := new(mock.MockProductRepo)
			tc.setupMock(mockProductRepo)

			svc := NewProductService(mockProductRepo)

			products, err := svc.GetAll(ctx)
			if err != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorContains)
			} else {
				require.NoError(t, err)
				require.Len(t, products, tc.expectedLen)
			}

			mockProductRepo.AssertExpectations(t)
		})
	}
}
