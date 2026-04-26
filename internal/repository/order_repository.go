package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var orderID string
	err = tx.QueryRow(
		`INSERT INTO orders(user_id, total_price, status, receiver_name, phone, province, district, logistic, logistic_branch)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		order.UserID, order.TotalPrice, order.Status,
		order.ReceiverName, order.Phone,
		order.Province, order.District,
		order.Logistic, order.LogisticBranch,
	).Scan(&orderID)
	if err != nil {
		return "", err
	}

	for _, item := range order.OrderItems {
		_, err = tx.Exec(
			`INSERT INTO order_items(order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`,
			orderID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			return "", err
		}
		
		// Deduct stock
		_, err = tx.Exec(
			`UPDATE products SET stock = stock - $1 WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return "", err
		}
	}

	return orderID, tx.Commit()
}

func (r *OrderRepository) GetOrderByID(id string) (*domain.Order, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, total_price, status, receiver_name, phone,
		        province, district, logistic, logistic_branch, created_at, updated_at
		 FROM orders WHERE id=$1`, id,
	)
	var o domain.Order
	err := row.Scan(
		&o.ID, &o.UserID, &o.TotalPrice, &o.Status,
		&o.ReceiverName, &o.Phone,
		&o.Province, &o.District,
		&o.Logistic, &o.LogisticBranch,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	items, err := r.GetOrderItems(id)
	if err != nil {
		return nil, err
	}
	o.OrderItems = items
	return &o, nil
}

func (r *OrderRepository) GetOrderItems(orderID string) ([]domain.OrderItem, error) {
	rows, err := r.db.Query(
		`SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id=$1`, orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID string) ([]domain.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, total_price, status, receiver_name, phone,
		        province, district, logistic, logistic_branch, created_at, updated_at
		 FROM orders WHERE user_id=$1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.TotalPrice, &o.Status,
			&o.ReceiverName, &o.Phone,
			&o.Province, &o.District,
			&o.Logistic, &o.LogisticBranch,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(id string, status string) error {
	res, err := r.db.Exec(
		`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, status, id,
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

func (r *OrderRepository) GetAllOrders() ([]domain.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, total_price, status, receiver_name, phone,
		        province, district, logistic, logistic_branch, created_at, updated_at
		 FROM orders ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.TotalPrice, &o.Status,
			&o.ReceiverName, &o.Phone,
			&o.Province, &o.District,
			&o.Logistic, &o.LogisticBranch,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
