package service

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrProductNotFound      = errors.New("product not found")
	ErrAlreadyPaused        = errors.New("subscription is already paused")
	ErrAlreadyCancelled     = errors.New("subscription is already cancelled")
	ErrAlreadyActive        = errors.New("subscription is already active")
	ErrInvalidState         = errors.New("can't pause at this state")
	ErrNoPendingPayment     = errors.New("no pending payment for this subscription")
	ErrFailedPayment        = errors.New("payment failed")
)
