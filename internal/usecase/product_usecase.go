package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type ProductUsecase struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (u *ProductUsecase) CreateProduct(sellerID, name, description string, price float64, stock int, status string, imageURLs []string, categoryIDs []string) (string, error) {
	product := &domain.Product{
		SellerID:    sellerID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Status:      status,
	}
	if product.Status == "" {
		product.Status = "active"
	}

	productID, err := u.productRepo.CreateProduct(product)
	if err != nil {
		return "", err
	}

	for _, imageURL := range imageURLs {
		_ = u.productRepo.AddProductImage(&domain.ProductImage{
			ProductID: productID,
			ImageURL:  imageURL,
		})
	}

	for _, categoryID := range categoryIDs {
		_ = u.productRepo.AddProductCategory(productID, categoryID)
	}

	return productID, nil
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