package usecase

import (
	"errors"
	"os"

	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
	"github.com/l0ng7h0r/golang/pkg/security"
)

type AuthUsecase struct {
	userRepo *repository.UserRepository
}

func NewAuthUsecase(userRepo *repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

// CreateUser

func (u *AuthUsecase) CreateUser(email, password string, roles []string) error {
	hashed, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email: email,
		Password: hashed,
		Roles: roles,
	}

	err = u.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID

func (u *AuthUsecase) GetUserByID(id string) (*domain.User, error) {
	return u.userRepo.GetUserByID(id)
}

// GetAllUsers

func (u *AuthUsecase) GetAllUsers() ([]domain.User, error) {
	return u.userRepo.GetAllUsers()
}

// DeleteUser

func (u *AuthUsecase) DeleteUser(id string) error {
	return u.userRepo.DeleteUser(id)
}

// RegisterUser

func (u *AuthUsecase) Register(email, password string) error {
	hashed, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: hashed,
		Roles:    []string{"user"},
	}

	err = u.userRepo.RegisterUser(user)
	if err != nil {
		return err
	}

	return nil
}	

// Login

func (u *AuthUsecase) Login(email, password string) (string, string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := security.VerifyPassword(user.Password, password); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := security.GenerateToken(user.ID, user.Roles, os.Getenv("JWT_ACCESS_EXPIRY"))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := security.GenerateToken(user.ID, user.Roles, os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		return "", "", err
	}

	err = u.userRepo.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Refresh

func (u *AuthUsecase) Refresh(token string) (string, string, error) {
	userID, err := u.userRepo.ValidateRefreshToken(token)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	_ = u.userRepo.DeleteRefreshToken(token)

	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	accessToken, err := security.GenerateToken(user.ID, user.Roles, os.Getenv("JWT_ACCESS_EXPIRY"))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := security.GenerateToken(user.ID, user.Roles, os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		return "", "", err
	}

	err = u.userRepo.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


