package domain

import "time"

type Seller struct {
	UserID      string    `json:"user_id"`
	StoreName   string    `json:"store_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}