package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type AddressUsecase struct {
	addressRepo *repository.AddressRepository
}

func NewAddressUsecase(addressRepo *repository.AddressRepository) *AddressUsecase {
	return &AddressUsecase{addressRepo: addressRepo}
}

func (u *AddressUsecase) CreateAddress(userID, name, phone, address, city, country, postalCode string, isDefault bool) (string, error) {
	return u.addressRepo.CreateAddress(&domain.Address{
		UserID:     userID,
		Name:       name,
		Phone:      phone,
		Address:    address,
		City:       city,
		Country:    country,
		PostalCode: postalCode,
		IsDefault:  isDefault,
	})
}

func (u *AddressUsecase) GetAddressByID(id string) (*domain.Address, error) {
	return u.addressRepo.GetAddressByID(id)
}

func (u *AddressUsecase) GetMyAddresses(userID string) ([]domain.Address, error) {
	return u.addressRepo.GetAddressesByUserID(userID)
}

func (u *AddressUsecase) UpdateAddress(id string, address *domain.Address) error {
	return u.addressRepo.UpdateAddress(id, address)
}

func (u *AddressUsecase) DeleteAddress(id, userID string) error {
	return u.addressRepo.DeleteAddress(id, userID)
}

func (u *AddressUsecase) SetDefaultAddress(id, userID string) error {
	return u.addressRepo.SetDefaultAddress(id, userID)
}
