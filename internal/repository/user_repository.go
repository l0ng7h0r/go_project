package repository

import (
	"database/sql"

	"github.com/l0ng7h0r/golang/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser

func (r *UserRepository) CreateUser(user *domain.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var userID string
	err = tx.QueryRow(`INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id`, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		_, err := tx.Exec(`INSERT INTO user_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name = $2`, userID, role)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

// CreateUserReturningID creates a user and returns the generated ID

func (r *UserRepository) CreateUserReturningID(user *domain.User) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	defer tx.Rollback()

	var userID string
	err = tx.QueryRow(`INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id`, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return "", err
	}

	for _, role := range user.Roles {
		_, err := tx.Exec(`INSERT INTO user_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name = $2`, userID, role)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return userID, nil
}

// GetUserByID

func (r *UserRepository) GetUserByID(id string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password FROM users WHERE id=$1`, id)
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`, u.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var role string
			if err := rows.Scan(&role); err == nil {
				u.Roles = append(u.Roles, role)
			}
		}
	}

	return &u, nil
}

// GetAllUsers

func (r *UserRepository) GetAllUsers() ([]domain.User, error) {
	rows, err := r.db.Query(`SELECT id, email, password, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		var u domain.User
		err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// DeleteUser

func (r *UserRepository) DeleteUser(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

// RegisterUser

func (r *UserRepository) RegisterUser(user *domain.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var userID string
	err = tx.QueryRow(`INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id`, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		_, err := tx.Exec(`INSERT INTO user_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name = $2`, userID, role)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

// FindByEmail

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password FROM users WHERE email=$1`, email)
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`, u.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var role string
			if err := rows.Scan(&role); err == nil {
				u.Roles = append(u.Roles, role)
			}
		}
	}

	return &u, nil
}

// FindByID

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password FROM users WHERE id=$1`, id)
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`, u.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var role string
			if err := rows.Scan(&role); err == nil {
				u.Roles = append(u.Roles, role)
			}
		}
	}

	return &u, nil
}

// SaveRefreshToken

func (r *UserRepository) SaveRefreshToken(userID, token string) error {
	_, err := r.db.Exec(`INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES ($1,$2,NOW() + INTERVAL '7 days')`, userID, token)
	return err
}

// ValidateRefreshToken

func (r *UserRepository) ValidateRefreshToken(token string) (string, error) {
	row := r.db.QueryRow(`SELECT user_id FROM refresh_tokens WHERE token=$1 AND expires_at > NOW()`, token)
	var userID string
	err := row.Scan(&userID)
	return userID, err
}

// DeleteRefreshToken

func (r *UserRepository) DeleteRefreshToken(token string) error {
	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE token=$1", token)
	return err
}