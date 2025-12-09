package service

import (
	"errors"
	"sync"
)

var ErrAlreadyProcessed = errors.New("transaction already processed")

type PaymentService struct {
	mu           sync.Mutex // Mutual Exclusion Lock to protect concurrent read/write
	transactions map[string]Transaction
}

// Init new paymentservice
func NewPaymentService() *PaymentService {
	return &PaymentService{
		transactions: make(map[string]Transaction),
	}
}

func (p *PaymentService) ProcessPayment(t Transaction) (Transaction, error) {
	// Lock resource, use defer to guarantee unlock
	p.mu.Lock()
	defer p.mu.Unlock()

	// idenpotency check
	if existingTx, ok := p.transactions[t.TransactionID]; ok {
		return existingTx, ErrAlreadyProcessed
	}

	t.Status = TransactionStatusSuccess
	p.transactions[t.TransactionID] = t

	return t, nil
}
