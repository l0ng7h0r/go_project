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

func (r *ProductRepository) CreateProduct(product *domain.Product) (string, error) {
	var productID string
	err := r.db.QueryRow(
		`INSERT INTO products(seller_id, name, description, price, stock, status)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		product.SellerID, product.Name, product.Description, product.Price, product.Stock, product.Status,
	).Scan(&productID)
	if err != nil {
		return "", err
	}
	return productID, nil
}

func (r *ProductRepository) AddProductImage(image *domain.ProductImage) error {
	_, err := r.db.Exec(
		`INSERT INTO product_images(product_id, image_url) VALUES ($1, $2)`,
		image.ProductID, image.ImageURL,
	)
	return err
}

func (r *ProductRepository) AddProductCategory(productID, categoryID string) error {
	_, err := r.db.Exec(
		`INSERT INTO product_categories(product_id, category_id) VALUES ($1, $2)`,
		productID, categoryID,
	)
	return err
}

func (r *ProductRepository) GetProductByID(id string) (*domain.Product, error) {
	row := r.db.QueryRow(
		`SELECT id, seller_id, name, description, price, stock, status, created_at, updated_at
		 FROM products WHERE id=$1`, id,
	)
	var p domain.Product
	err := row.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Load images
	imgRows, err := r.db.Query(`SELECT id, product_id, image_url, created_at FROM product_images WHERE product_id=$1`, id)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img domain.ProductImage
			if err := imgRows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &img.CreatedAt); err == nil {
				p.Images = append(p.Images, img)
			}
		}
	}

	return &p, nil
}

func (r *ProductRepository) GetAllProducts() ([]domain.Product, error) {
	rows, err := r.db.Query(
		`SELECT id, seller_id, name, description, price, stock, status, created_at, updated_at FROM products`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetProductsBySeller(sellerID string) ([]domain.Product, error) {
	rows, err := r.db.Query(
		`SELECT id, seller_id, name, description, price, stock, status, created_at, updated_at
		 FROM products WHERE seller_id=$1`, sellerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Status, &p.CreatedAt, &p.UpdatedAt)
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
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(id string, product *domain.Product) error {
	res, err := r.db.Exec(
		`UPDATE products SET name=$1, description=$2, price=$3, stock=$4, status=$5, updated_at=NOW()
		 WHERE id=$6 AND seller_id=$7`,
		product.Name, product.Description, product.Price, product.Stock, product.Status, id, product.SellerID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}