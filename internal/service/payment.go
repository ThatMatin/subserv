package service

import (
	"fmt"
	"math/rand"
)

type PaymentProcessor interface {
	Charge(req PaymentRequest) (*PaymentResult, error)
}

type PaymentRequest struct {
	UserID    uint
	ProductID uint
	Amount    int
}

type PaymentResult struct {
	Success bool
	TxID    string
	Error   string
}

type dummyPaymentProcessor struct{}

func NewDummyPaymentProcessor() PaymentProcessor {
	return &dummyPaymentProcessor{}
}

func (p *dummyPaymentProcessor) Charge(req PaymentRequest) (*PaymentResult, error) {
	// simulate 95% success rate
	// returned error will be used for internal errors but payment errors checked seperately
	if rand.Float64() < 0.95 {
		return &PaymentResult{
			Success: true,
			TxID:    fmt.Sprintf("tx-%d", rand.Intn(1000000)),
		}, nil
	}

	return &PaymentResult{
		Success: false,
		Error:   "payment failed",
	}, nil
}
