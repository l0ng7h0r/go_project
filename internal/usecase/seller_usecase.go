package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
	"github.com/l0ng7h0r/golang/pkg/security"
)

type SellerUsecase struct {
	sellerRepo *repository.SellerRepository
	userRepo   *repository.UserRepository
}

func NewSellerUsecase(sellerRepo *repository.SellerRepository, userRepo *repository.UserRepository) *SellerUsecase {
	return &SellerUsecase{sellerRepo: sellerRepo, userRepo: userRepo}
}

func (u *SellerUsecase) CreateSeller(email, password string, roles []string, seller *domain.Seller) error {
	hashed, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: hashed,
		Roles:    roles,
	}

	userID, err := u.userRepo.CreateUserReturningID(user)
	if err != nil {
		return err
	}

	seller.UserID = userID
	return u.sellerRepo.CreateSeller(seller)
}

func (u *SellerUsecase) GetSellerByID(id string) (*domain.Seller, error) {
	return u.sellerRepo.GetSellerByID(id)
}

func (u *SellerUsecase) GetAllSellers() ([]domain.Seller, error) {
	return u.sellerRepo.GetAllSellers()
}

func (u *SellerUsecase) DeleteSeller(id string) error {
	// First get the seller to find the linked user_id
	seller, err := u.sellerRepo.GetSellerByID(id)
	if err != nil {
		return err
	}

	// Delete seller first (child record)
	if err := u.sellerRepo.DeleteSeller(id); err != nil {
		return err
	}

	// Then delete the linked user (parent record)
	return u.userRepo.DeleteUser(seller.UserID)
}

func (u *SellerUsecase) UpdateSeller(id string, seller *domain.Seller) error {
	return u.sellerRepo.UpdateSeller(id, seller)
}