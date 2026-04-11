package repository

import(
	"database/sql"
	"github.com/l0ng7h0r/golang/internal/domain"
)

type SellerRepository struct {
	db *sql.DB
}

func NewSellerRepository(db *sql.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

func (r *SellerRepository) CreateSeller(seller *domain.Seller) error {
	_, err := r.db.Exec(`INSERT INTO sellers(user_id, store_name, description) VALUES ($1, $2, $3)`, seller.UserID, seller.StoreName, seller.Description)
	if err != nil {
		return err
	}
	return nil
}

func (r *SellerRepository) GetSellerByID(id string) (*domain.Seller, error) {
	row := r.db.QueryRow(`SELECT user_id, store_name, description, created_at, updated_at FROM sellers WHERE user_id=$1`, id)
	var s domain.Seller
	err := row.Scan(&s.UserID, &s.StoreName, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SellerRepository) GetAllSellers() ([]domain.Seller, error) {
	rows, err := r.db.Query(`SELECT user_id, store_name, description, created_at, updated_at FROM sellers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sellers []domain.Seller
	for rows.Next() {
		var s domain.Seller
		err := rows.Scan(&s.UserID, &s.StoreName, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}
	return sellers, nil
}

func (r *SellerRepository) DeleteSeller(id string) error {
	_, err := r.db.Exec(`DELETE FROM sellers WHERE user_id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *SellerRepository) UpdateSeller(id string, seller *domain.Seller) error {
	_, err := r.db.Exec(`UPDATE sellers SET user_id=$1, store_name=$2, description=$3 WHERE user_id=$4`, seller.UserID, seller.StoreName, seller.Description, id)
	if err != nil {
		return err
	}
	return nil
}