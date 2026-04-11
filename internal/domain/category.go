package domain

type Category struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"` // pointer to allow null for root categories
}
