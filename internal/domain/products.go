package domain

import "time"

type Product struct {
	ID          string         `json:"id"`
	SellerID    string         `json:"seller_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	Status      string         `json:"status"`
	Images      []ProductImage `json:"images"`
	Categories  []string       `json:"categories"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}