package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetOrCreateCart returns existing cart or creates a new one for the user
func (r *CartRepository) GetOrCreateCart(userID string) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.QueryRow(
		`SELECT id, user_id, created_at FROM carts WHERE user_id=$1`, userID,
	).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt)

	if err == sql.ErrNoRows {
		// Create new cart
		err = r.db.QueryRow(
			`INSERT INTO carts(user_id) VALUES ($1) RETURNING id, user_id, created_at`, userID,
		).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Load items
	items, err := r.GetCartItems(cart.ID)
	if err != nil {
		return nil, err
	}
	cart.Items = items
	return &cart, nil
}

func (r *CartRepository) GetCartItems(cartID string) ([]domain.CartItem, error) {
	rows, err := r.db.Query(
		`SELECT cart_id, product_id, quantity FROM cart_items WHERE cart_id=$1`, cartID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.CartItem
	for rows.Next() {
		var item domain.CartItem
		if err := rows.Scan(&item.CartID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *CartRepository) AddOrUpdateCartItem(cartID, productID string, quantity int) error {
	_, err := r.db.Exec(
		`INSERT INTO cart_items(cart_id, product_id, quantity)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + $3`,
		cartID, productID, quantity,
	)
	return err
}

func (r *CartRepository) UpdateCartItemQuantity(cartID, productID string, quantity int) error {
	if quantity <= 0 {
		return r.RemoveCartItem(cartID, productID)
	}
	res, err := r.db.Exec(
		`UPDATE cart_items SET quantity=$1 WHERE cart_id=$2 AND product_id=$3`,
		quantity, cartID, productID,
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

func (r *CartRepository) RemoveCartItem(cartID, productID string) error {
	_, err := r.db.Exec(
		`DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2`, cartID, productID,
	)
	return err
}

func (r *CartRepository) ClearCart(cartID string) error {
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE cart_id=$1`, cartID)
	return err
}
