package domain

import "time"

type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	AddressID  string      `json:"address_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"` // pending, confirmed, shipped, delivered, cancelled
	OrderItems []OrderItem `json:"order_items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // snapshot price at order time
	Subtotal  float64 `json:"subtotal"`
}
