package domain

import "time"

type Address struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Label       string    `json:"label"` // e.g. Home, Office
	AddressLine string    `json:"address_line"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
}
