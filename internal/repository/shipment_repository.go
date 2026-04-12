package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) CreateShipment(shipment *domain.Shipment) (string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO shipments(order_id, provider, tracking_number, status)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		shipment.OrderID, shipment.Provider, shipment.TrackingNumber, shipment.Status,
	).Scan(&id)
	return id, err
}

func (r *ShipmentRepository) GetShipmentByOrderID(orderID string) (*domain.Shipment, error) {
	row := r.db.QueryRow(
		`SELECT id, order_id, provider, tracking_number, status, shipped_at, delivered_at, created_at
		 FROM shipments WHERE order_id=$1`, orderID,
	)
	var s domain.Shipment
	err := row.Scan(
		&s.ID, &s.OrderID, &s.Provider, &s.TrackingNumber, &s.Status,
		&s.ShippedAt, &s.DeliveredAt, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShipmentRepository) UpdateShipmentStatus(id string, status string) error {
	var query string
	switch status {
	case "shipped":
		query = `UPDATE shipments SET status=$1, shipped_at=NOW() WHERE id=$2`
	case "delivered":
		query = `UPDATE shipments SET status=$1, delivered_at=NOW() WHERE id=$2`
	default:
		query = `UPDATE shipments SET status=$1 WHERE id=$2`
	}
	res, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ShipmentRepository) UpdateTrackingNumber(id string, provider string, trackingNumber string) error {
	res, err := r.db.Exec(
		`UPDATE shipments SET provider=$1, tracking_number=$2 WHERE id=$3`,
		provider, trackingNumber, id,
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
