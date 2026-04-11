package domain

import "time"

type UserProfile struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
