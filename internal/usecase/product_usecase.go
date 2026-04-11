package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)


type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) CreateProduct(sellerID string, name string, description string, price float64, stock int) error {
	product := &domain.Product{
		SellerID: sellerID,
		Name: name,
		Description: description,
		Price: price,
		Stock: stock,
	}
	return u.productRepo.CreateProduct(product)
}

func (u *ProductUsecase) GetProductByID(id string) (*domain.Product, error) {
	return u.productRepo.GetProductByID(id)
}

func (u *ProductUsecase) GetAllProducts() ([]domain.Product, error) {
	return u.productRepo.GetAllProducts()
}

func (u *ProductUsecase) GetProductsBySeller(sellerID string) ([]domain.Product, error) {
	return u.productRepo.GetProductsBySeller(sellerID)
}

func (u *ProductUsecase) DeleteProduct(id string, sellerID string) error {
	return u.productRepo.DeleteProduct(id, sellerID)
}

func (u *ProductUsecase) UpdateProduct(id string, product *domain.Product) error {
	return u.productRepo.UpdateProduct(id, product)
}