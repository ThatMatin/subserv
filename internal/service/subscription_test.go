package service

import (
	"context"
	"errors"
	"testing"
	"time"

	mocklib "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thatmatin/subserv/internal/mock"
	"github.com/thatmatin/subserv/internal/model"
	"gorm.io/gorm"
)

var fixedTime = time.Date(2020, time.May, 0, 0, 0, 0, 0, time.UTC)

func TestFetchSubscriptionInfo(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name             string
		inputID          uint
		expectedUserID   uint
		expectedErr      error
		errorContains    string
		expectedStart    time.Time
		expectedState    model.State
		expectedDuration time.Duration
		expectedPrice    int
		setupMock        func(repo *mock.MockSubscriptionRepo)
	}{
		{
			name:             "existing subscription",
			inputID:          1,
			expectedUserID:   1,
			expectedErr:      nil,
			expectedState:    model.Pending,
			expectedStart:    fixedTime,
			expectedDuration: time.Hour * 24 * 30,
			expectedPrice:    1000,
			setupMock: func(repo *mock.MockSubscriptionRepo) {
				subscription := &model.Subscription{
					Model:     gorm.Model{ID: 1},
					UserID:    uint(1),
					Start:     fixedTime,
					End:       fixedTime.Add(time.Hour * 24 * 30),
					State:     model.Pending,
					PriceCent: 1000,
					TaxRate:   10,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
		{
			name:        "non-existent subscription",
			inputID:     2,
			expectedErr: ErrSubscriptionNotFound,
			setupMock: func(repo *mock.MockSubscriptionRepo) {
				repo.On("GetByID", ctx, uint(2)).Return((*model.Subscription)(nil), ErrSubscriptionNotFound)
			},
		},
		{
			name:          "repository error",
			inputID:       2,
			expectedErr:   errors.New(""),
			errorContains: "failed to fetch subscription",
			setupMock: func(repo *mock.MockSubscriptionRepo) {
				repo.On("GetByID", ctx, uint(2)).Return((*model.Subscription)(nil), errors.New("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := new(mock.MockSubscriptionRepo)
			u := new(mock.MockUserRepo)
			p := new(mock.MockProductRepo)
			tc.setupMock(s)

			svc := NewSubscriptionService(s, &productService{p}, &userService{u}, &dummyPaymentProcessor{})

			subscription, err := svc.Get(ctx, tc.inputID)
			if tc.expectedErr != nil {
				require.Error(t, err)
				if tc.errorContains == "" {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.ErrorContains(t, err, tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUserID, subscription.UserID)
				require.Equal(t, tc.expectedStart, subscription.Start)
				require.Equal(t, tc.expectedState, subscription.State)
				require.Equal(t, tc.expectedDuration, subscription.End.Sub(subscription.Start))
				require.Equal(t, tc.expectedPrice, subscription.PriceCent)
			}

			s.AssertExpectations(t)
			u.AssertExpectations(t)
			p.AssertExpectations(t)
		})
	}
}

func TestCreateSubscription(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name          string
		productID     uint
		userID        uint
		expectedState model.State
		expectedErr   error
		errorContains string
		setupMock     func(subsRepo *mock.MockSubscriptionRepo, prodRepo *mock.MockProductRepo, userRepo *mock.MockUserRepo)
	}{
		{
			name:          "successful",
			productID:     1,
			userID:        2,
			expectedErr:   nil,
			expectedState: model.Pending,
			setupMock: func(subsRepo *mock.MockSubscriptionRepo, prodRepo *mock.MockProductRepo, userRepo *mock.MockUserRepo) {
				product := &model.Product{
					Model:    gorm.Model{ID: 1},
					Name:     "flowmotion",
					Duration: time.Hour * 24 * 30,
					Price:    10000,
					TaxRate:  10,
				}
				userRepo.On("Exists", ctx, uint(2)).Return(true, nil)
				prodRepo.On("GetByID", ctx, uint(1)).Return(product, nil)
				subsRepo.On("Create", ctx, mocklib.MatchedBy(func(sub *model.Subscription) bool {
					return sub.UserID == uint(2) &&
						sub.ProductID == uint(1) &&
						sub.State == model.Pending &&
						sub.PriceCent == product.Price &&
						sub.TaxRate == product.TaxRate &&
						!sub.Start.IsZero() &&
						!sub.End.IsZero()
				})).Return(nil)
			},
		},
		{
			name:        "user not found",
			userID:      2,
			expectedErr: ErrUserNotFound,
			setupMock: func(subsRepo *mock.MockSubscriptionRepo, prodRepo *mock.MockProductRepo, userRepo *mock.MockUserRepo) {
				userRepo.On("Exists", ctx, uint(2)).Return(false, nil)
			},
		},
		{
			name:        "product not found",
			userID:      2,
			productID:   2,
			expectedErr: ErrProductNotFound,
			setupMock: func(subsRepo *mock.MockSubscriptionRepo, prodRepo *mock.MockProductRepo, userRepo *mock.MockUserRepo) {
				userRepo.On("Exists", ctx, uint(2)).Return(true, nil)
				prodRepo.On("GetByID", ctx, uint(2)).Return((*model.Product)(nil), ErrProductNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := new(mock.MockSubscriptionRepo)
			p := new(mock.MockProductRepo)
			u := new(mock.MockUserRepo)
			tc.setupMock(s, p, u)

			svc := NewSubscriptionService(s, &productService{p}, &userService{u}, &dummyPaymentProcessor{})

			subscription, err := svc.Create(ctx, tc.productID, tc.userID)
			if tc.expectedErr != nil {
				require.Error(t, err)
				if tc.errorContains == "" {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.ErrorContains(t, err, tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.userID, subscription.UserID)
				require.Equal(t, tc.expectedState, subscription.State)
			}

			s.AssertExpectations(t)
			u.AssertExpectations(t)
			p.AssertExpectations(t)
		})
	}
}

func TestPauseSubscription(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		expectedErr error
		status      model.State
		setupMock   func(repo *mock.MockSubscriptionRepo, status model.State)
	}{
		{
			name:        "pause active subscription",
			expectedErr: nil,
			status:      model.Active,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				start := time.Now()
				end := start.Add(time.Hour * 24 * 30)
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
					Start: start,
					End:   end,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
				repo.On("Save", ctx, mocklib.MatchedBy(func(subs *model.Subscription) bool { return subs.State == model.Paused })).Return(nil)
			},
		},
		{
			name:        "pause paused subscription",
			expectedErr: ErrAlreadyPaused,
			status:      model.Paused,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
		{
			name:        "pause cancelled subscription",
			expectedErr: ErrInvalidState,
			status:      model.Cancelled,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
		{
			name:        "pausing expired subscription",
			expectedErr: ErrInvalidState,
			status:      model.Expired,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mock.MockSubscriptionRepo)
			tc.setupMock(repo, tc.status)
			svc := NewSubscriptionService(repo, &productService{}, &userService{}, &dummyPaymentProcessor{})

			if err := svc.Pause(ctx, 1); err != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		expectedErr error
		status      model.State
		setupMock   func(repo *mock.MockSubscriptionRepo, state model.State)
	}{
		{
			name:        "cancel active subscription",
			expectedErr: nil,
			status:      model.Active,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				start := time.Now()
				end := start.Add(time.Hour * 24 * 30)
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
					Start: start,
					End:   end,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
				repo.On("Save", ctx, mocklib.MatchedBy(func(subs *model.Subscription) bool { return subs.State == model.Cancelled })).Return(nil)
			},
		},
		{
			name:        "cancel pending subscription",
			expectedErr: nil,
			status:      model.Pending,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				start := time.Now()
				end := start.Add(time.Hour * 24 * 30)
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
					Start: start,
					End:   end,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
				repo.On("Save", ctx, mocklib.MatchedBy(func(subs *model.Subscription) bool { return subs.State == model.Cancelled })).Return(nil)
			},
		},

		{
			name:        "cancel paused subscription",
			expectedErr: nil,
			status:      model.Paused,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				start := time.Now()
				end := start.Add(time.Hour * 24 * 30)
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
					Start: start,
					End:   end,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
				repo.On("Save", ctx, mocklib.MatchedBy(func(subs *model.Subscription) bool { return subs.State == model.Cancelled })).Return(nil)
			},
		},
		{
			name:        "cancel cancelled subscription",
			expectedErr: ErrAlreadyCancelled,
			status:      model.Cancelled,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
		{
			name:        "cancel expired subscription",
			expectedErr: ErrInvalidState,
			status:      model.Expired,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mock.MockSubscriptionRepo)
			tc.setupMock(repo, tc.status)
			svc := NewSubscriptionService(repo, &productService{}, &userService{}, &dummyPaymentProcessor{})

			if err := svc.Cancel(ctx, 1); err != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestUnpauseSubscription(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		expectedErr error
		state       model.State
		setupMock   func(repo *mock.MockSubscriptionRepo, state model.State)
	}{
		{
			name:        "unpause paused subscription",
			expectedErr: nil,
			state:       model.Paused,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				start := time.Now().Add(-time.Hour * 48)
				end := start.Add(time.Hour * 24 * 2)
				pauseAt := start.Add(time.Hour * 24)
				subscription := &model.Subscription{
					Model:    gorm.Model{ID: 1},
					State:    state,
					Start:    start,
					End:      end,
					PausedAt: &pauseAt,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
				repo.On("Save", ctx, mocklib.MatchedBy(func(subs *model.Subscription) bool { return subs.State == model.Active })).Return(nil)
			},
		},
		{
			name:        "unpause active subscription",
			expectedErr: ErrAlreadyActive,
			state:       model.Active,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
		{
			name:        "unpausing cancelled subscription",
			expectedErr: ErrInvalidState,
			state:       model.Cancelled,
			setupMock: func(repo *mock.MockSubscriptionRepo, state model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: state,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
		{
			name:        "unpausing expired subscription",
			expectedErr: ErrInvalidState,
			state:       model.Expired,
			setupMock: func(repo *mock.MockSubscriptionRepo, status model.State) {
				subscription := &model.Subscription{
					Model: gorm.Model{ID: 1},
					State: status,
				}
				repo.On("GetByID", ctx, uint(1)).Return(subscription, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mock.MockSubscriptionRepo)
			tc.setupMock(repo, tc.state)
			svc := NewSubscriptionService(repo, &productService{}, &userService{}, &dummyPaymentProcessor{})

			if err := svc.Unpause(ctx, 1); err != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
