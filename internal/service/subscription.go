package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thatmatin/subserv/internal/model"
	"github.com/thatmatin/subserv/internal/repo"
	"github.com/thatmatin/subserv/internal/utils"
	"gorm.io/gorm"
)

var UTCLocation *time.Location

type SubscriptionService interface {
	Get(ctx context.Context, ID uint) (*model.Subscription, error)
	Create(ctx context.Context, productID uint, userID uint) (*model.Subscription, error)
	Purchase(ctx context.Context, ID uint) error
	Pause(ctx context.Context, ID uint) error
	Unpause(ctx context.Context, ID uint) error
	Cancel(ctx context.Context, ID uint) error
}

type subscriptionService struct {
	subsRepo         repo.SubscriptionRepository
	productService   ProductService
	userService      UserService
	paymentProcessor PaymentProcessor
}

func NewSubscriptionService(
	subsRepo repo.SubscriptionRepository,
	prodSvc ProductService,
	userSvc UserService,
	paySvc PaymentProcessor,
) SubscriptionService {
	return &subscriptionService{
		subsRepo:         subsRepo,
		productService:   prodSvc,
		userService:      userSvc,
		paymentProcessor: paySvc,
	}
}

func (s *subscriptionService) Get(ctx context.Context, ID uint) (*model.Subscription, error) {
	subscription, err := s.subsRepo.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	return subscription, nil
}

func (s *subscriptionService) Create(ctx context.Context, productID uint, userID uint) (*model.Subscription, error) {
	exists, err := s.userService.Exists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	product, err := s.productService.Get(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("couldn't fetch product: %w", err)
	}

	now := time.Now().In(UTCLocation)
	subscription := &model.Subscription{
		UserID:    userID,
		ProductID: productID,
		Start:     now,
		End:       now.Add(product.Duration * 1e9), // Duration is in seconds, convert to nanoseconds
		State:     model.Pending,
		PriceCent: product.Price,
		TaxRate:   product.TaxRate,
	}

	if err := s.subsRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("couldn't create subscription: %w", err)
	}

	return subscription, nil
}

func (s *subscriptionService) Purchase(ctx context.Context, ID uint) error {
	subscription, err := s.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubscriptionNotFound
		}

		return fmt.Errorf("couldn't pause subscription: %w", err)
	}

	if subscription.State != model.Pending {
		return ErrNoPendingPayment
	}

	payResult, err := s.paymentProcessor.Charge(PaymentRequest{
		UserID:    subscription.UserID,
		ProductID: subscription.ID,
		Amount:    utils.CalculateFinalAmount(subscription.PriceCent, subscription.TaxRate),
	})
	if err != nil {
		return fmt.Errorf("an error occured in payment: %w", err)
	}
	if !payResult.Success {
		return ErrFailedPayment
	}

	// idempotency and log of successful payment must be implemented in real payment implementation
	subscription.State = model.Active
	duration := subscription.End.Sub(subscription.Start)
	subscription.Start = time.Now().In(UTCLocation)
	subscription.End = subscription.Start.Add(duration)
	if err := s.subsRepo.Save(ctx, subscription); err != nil {
		return fmt.Errorf("couldn't save successful payment [Transaction ID %s] : %w", payResult.TxID, err)
	}

	return nil
}

func (s *subscriptionService) Pause(ctx context.Context, ID uint) error {
	subscription, err := s.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubscriptionNotFound
		}

		return fmt.Errorf("couldn't pause subscription: %w", err)
	}

	switch subscription.State {
	case model.Active:
		if time.Now().In(UTCLocation).Before(subscription.End.In(UTCLocation)) {
			subscription.State = model.Paused
			now := time.Now().In(UTCLocation)
			subscription.PausedAt = &now

			if err := s.subsRepo.Save(ctx, subscription); err != nil {
				return fmt.Errorf("couldn't pause subscription: %w", err)
			}

			return nil
		}
		return ErrAlreadyExpired
	case model.Paused:
		return ErrAlreadyPaused
	default:
		return ErrInvalidState
	}
}

func (s *subscriptionService) Cancel(ctx context.Context, ID uint) error {
	subscription, err := s.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubscriptionNotFound
		}

		return fmt.Errorf("couldn't cancel subscription: %w", err)
	}

	switch subscription.State {
	case model.Active, model.Pending, model.Paused:
		if time.Now().In(UTCLocation).Before(subscription.End.In(UTCLocation)) {
			subscription.State = model.Cancelled
			now := time.Now().In(UTCLocation)
			subscription.End = now
			if err := s.subsRepo.Save(ctx, subscription); err != nil {
				return fmt.Errorf("couldn't cancel subscription: %w", err)
			}

			return nil
		}
		return ErrAlreadyExpired
	case model.Cancelled:
		return ErrAlreadyCancelled
	default:
		return ErrInvalidState
	}
}

func (s *subscriptionService) Unpause(ctx context.Context, ID uint) error {
	subscription, err := s.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubscriptionNotFound
		}

		return fmt.Errorf("couldn't unpause subscription: %w", err)
	}

	switch subscription.State {
	case model.Paused:
		subscription.State = model.Active
		now := time.Now().In(UTCLocation)
		// TODO: Handle case where PausedAt is nil
		pausedDuration := now.Sub(*subscription.PausedAt)
		subscription.End = subscription.End.Add(pausedDuration)

		if err := s.subsRepo.Save(ctx, subscription); err != nil {
			return fmt.Errorf("couldn't cancel subscription: %w", err)
		}

		return nil
	case model.Active:
		return ErrAlreadyActive
	default:
		return ErrInvalidState
	}
}

func init() {
	var err error
	UTCLocation, err = time.LoadLocation("UTC")
	if err != nil {
		panic("Failed to load UTC location: " + err.Error())
	}
}
