package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category *domain.Category) (string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO categories(name, parent_id) VALUES ($1, $2) RETURNING id`,
		category.Name, category.ParentID,
	).Scan(&id)
	return id, err
}

func (r *CategoryRepository) GetCategoryByID(id string) (*domain.Category, error) {
	row := r.db.QueryRow(
		`SELECT id, name, parent_id, created_at FROM categories WHERE id=$1`, id,
	)
	var c domain.Category
	err := row.Scan(&c.ID, &c.Name, &c.ParentID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) GetAllCategories() ([]domain.Category, error) {
	rows, err := r.db.Query(
		`SELECT id, name, parent_id, created_at FROM categories ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(id string, category *domain.Category) error {
	res, err := r.db.Exec(
		`UPDATE categories SET name=$1, parent_id=$2 WHERE id=$3`,
		category.Name, category.ParentID, id,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id=$1`, id)
	return err
}
