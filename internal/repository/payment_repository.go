package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(payment *domain.Payment) (string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO payments(order_id, method, status, amount, transaction_id, payment_url, paid_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		payment.OrderID, payment.Method, payment.Status, payment.Amount,
		payment.TransactionID, payment.PaymentURL, payment.PaidAt,
	).Scan(&id)
	return id, err
}

func (r *PaymentRepository) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	row := r.db.QueryRow(
		`SELECT id, order_id, method, status, amount, transaction_id, payment_url, paid_at, created_at
		 FROM payments WHERE order_id=$1`, orderID,
	)
	var p domain.Payment
	err := row.Scan(
		&p.ID, &p.OrderID, &p.Method, &p.Status, &p.Amount,
		&p.TransactionID, &p.PaymentURL, &p.PaidAt, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) GetPaymentByID(id string) (*domain.Payment, error) {
	row := r.db.QueryRow(
		`SELECT id, order_id, method, status, amount, transaction_id, payment_url, paid_at, created_at
		 FROM payments WHERE id=$1`, id,
	)
	var p domain.Payment
	err := row.Scan(
		&p.ID, &p.OrderID, &p.Method, &p.Status, &p.Amount,
		&p.TransactionID, &p.PaymentURL, &p.PaidAt, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(id string, status string, transactionID string) error {
	res, err := r.db.Exec(
		`UPDATE payments SET status=$1, transaction_id=$2, paid_at=CASE WHEN $1='paid' THEN NOW() ELSE paid_at END
		 WHERE id=$3`,
		status, transactionID, id,
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
