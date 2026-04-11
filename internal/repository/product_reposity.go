package repository

import (
	"database/sql"
	"github.com/l0ng7h0r/golang/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product *domain.Product) error {
	_, err := r.db.Exec(`INSERT INTO products(seller_id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5)`, product.SellerID, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetProductByID(id string) (*domain.Product, error) {
	row := r.db.QueryRow(`SELECT id, seller_id, name, description, price, stock, created_at, updated_at FROM products WHERE id=$1`, id)
	var p domain.Product
	err := row.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) GetAllProducts() ([]domain.Product, error) {
	rows, err := r.db.Query(`SELECT id, seller_id, name, description, price, stock, created_at, updated_at FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetProductsBySeller(sellerID string) ([]domain.Product, error) {
	rows, err := r.db.Query(`SELECT id, seller_id, name, description, price, stock, created_at, updated_at FROM products WHERE seller_id=$1`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) DeleteProduct(id string, sellerID string) error {
	res, err := r.db.Exec(`DELETE FROM products WHERE id=$1 AND seller_id=$2`, id, sellerID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Or custom error
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(id string, product *domain.Product) error {
	res, err := r.db.Exec(`UPDATE products SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5 AND seller_id=$6`, product.Name, product.Description, product.Price, product.Stock, id, product.SellerID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Or custom error
	}
	return nil
}