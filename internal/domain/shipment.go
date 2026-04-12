package domain

import "time"

type Shipment struct {
	ID             string     `json:"id"`
	OrderID        string     `json:"order_id"`
	Provider       string     `json:"provider"`
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"` // pending, shipped, delivered
	ShippedAt      *time.Time `json:"shipped_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`
	CreatedAt      time.Time  `json:"created_at"`
}
