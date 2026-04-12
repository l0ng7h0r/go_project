package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type AddressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) CreateAddress(address *domain.Address) (string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO addresses(user_id, name, phone, address, city, country, postal_code, is_default)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		address.UserID, address.Name, address.Phone, address.Address,
		address.City, address.Country, address.PostalCode, address.IsDefault,
	).Scan(&id)
	return id, err
}

func (r *AddressRepository) GetAddressByID(id string) (*domain.Address, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, name, phone, address, city, country, postal_code, is_default, created_at
		 FROM addresses WHERE id=$1`, id,
	)
	var a domain.Address
	err := row.Scan(
		&a.ID, &a.UserID, &a.Name, &a.Phone, &a.Address,
		&a.City, &a.Country, &a.PostalCode, &a.IsDefault, &a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AddressRepository) GetAddressesByUserID(userID string) ([]domain.Address, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, phone, address, city, country, postal_code, is_default, created_at
		 FROM addresses WHERE user_id=$1 ORDER BY is_default DESC, created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var addresses []domain.Address
	for rows.Next() {
		var a domain.Address
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Name, &a.Phone, &a.Address,
			&a.City, &a.Country, &a.PostalCode, &a.IsDefault, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (r *AddressRepository) UpdateAddress(id string, address *domain.Address) error {
	res, err := r.db.Exec(
		`UPDATE addresses SET name=$1, phone=$2, address=$3, city=$4, country=$5, postal_code=$6, is_default=$7
		 WHERE id=$8 AND user_id=$9`,
		address.Name, address.Phone, address.Address, address.City,
		address.Country, address.PostalCode, address.IsDefault, id, address.UserID,
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

func (r *AddressRepository) DeleteAddress(id string, userID string) error {
	res, err := r.db.Exec(`DELETE FROM addresses WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *AddressRepository) SetDefaultAddress(id string, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset all defaults for this user
	_, err = tx.Exec(`UPDATE addresses SET is_default=false WHERE user_id=$1`, userID)
	if err != nil {
		return err
	}

	// Set the specified address as default
	res, err := tx.Exec(`UPDATE addresses SET is_default=true WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}
