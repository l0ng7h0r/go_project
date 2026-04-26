package domain

import "time"

type Payment struct {
	ID            string     `json:"id"`
	OrderID       string     `json:"order_id"`
	Method        string     `json:"method"`         // phajay
	Status        string     `json:"status"`         // pending, paid, failed, refunded
	Amount        float64    `json:"amount"`
	TransactionID string     `json:"transaction_id"` // from Phajay webhook
	PaymentURL    string     `json:"payment_url"`    // Phajay payment link
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
