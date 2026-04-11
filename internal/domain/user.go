package domain

import "time"

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	Roles []string `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}