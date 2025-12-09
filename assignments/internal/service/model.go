package service

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusFailed  TransactionStatus = "FAILED"
)

type Transaction struct {
	UserID        string            `json:"user_id"`
	Amount        float64           `json:"amount"`
	TransactionID string            `json:"transaction_id"`
	Status        TransactionStatus `json:"status"`
}
