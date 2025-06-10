package service

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrProductNotFound      = errors.New("product not found")
	ErrNoPendingPayment     = errors.New("no pending payment for this subscription")
	ErrFailedPayment        = errors.New("payment failed")
	ErrUnauthorizedAccess   = errors.New("unauthorized access on subscription")

	ErrInvalidState     = errors.New("forbidden action at this state")
	ErrAlreadyPaused    = fmt.Errorf("subscription is already paused: %w", ErrInvalidState)
	ErrAlreadyCancelled = fmt.Errorf("subscription is already cancelled: %w", ErrInvalidState)
	ErrAlreadyActive    = fmt.Errorf("subscription is already active: %w", ErrInvalidState)
	ErrAlreadyExpired   = fmt.Errorf("subscription is already expired: %w", ErrInvalidState)
)
