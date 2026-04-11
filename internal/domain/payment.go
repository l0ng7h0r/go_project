package domain

import "time"

type Payment struct {
	ID             string    `json:"id"`
	OrderID        string    `json:"order_id"`
	Method         string    `json:"method"` // credit_card, bank_transfer, promptpay
	Status         string    `json:"status"` // pending, paid, failed, refunded
	Amount         float64   `json:"amount"`
	TransactionRef string    `json:"transaction_ref"`
	CreatedAt      time.Time `json:"created_at"`
}
