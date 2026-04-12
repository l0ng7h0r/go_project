package usecase

import (
	"errors"

	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type CartUsecase struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartUsecase(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartUsecase {
	return &CartUsecase{cartRepo: cartRepo, productRepo: productRepo}
}

func (u *CartUsecase) GetCart(userID string) (*domain.Cart, error) {
	return u.cartRepo.GetOrCreateCart(userID)
}

func (u *CartUsecase) AddItem(userID, productID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	// Verify product exists
	_, err := u.productRepo.GetProductByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}
	return u.cartRepo.AddOrUpdateCartItem(cart.ID, productID, quantity)
}

func (u *CartUsecase) UpdateItem(userID, productID string, quantity int) error {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}
	return u.cartRepo.UpdateCartItemQuantity(cart.ID, productID, quantity)
}

func (u *CartUsecase) RemoveItem(userID, productID string) error {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}
	return u.cartRepo.RemoveCartItem(cart.ID, productID)
}

func (u *CartUsecase) ClearCart(userID string) error {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}
	return u.cartRepo.ClearCart(cart.ID)
}
