package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type CategoryUsecase struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo *repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{categoryRepo: categoryRepo}
}

func (u *CategoryUsecase) CreateCategory(name string, parentID *string) (string, error) {
	return u.categoryRepo.CreateCategory(&domain.Category{
		Name:     name,
		ParentID: parentID,
	})
}

func (u *CategoryUsecase) GetCategoryByID(id string) (*domain.Category, error) {
	return u.categoryRepo.GetCategoryByID(id)
}

func (u *CategoryUsecase) GetAllCategories() ([]domain.Category, error) {
	return u.categoryRepo.GetAllCategories()
}

func (u *CategoryUsecase) UpdateCategory(id string, name string, parentID *string) error {
	return u.categoryRepo.UpdateCategory(id, &domain.Category{
		Name:     name,
		ParentID: parentID,
	})
}

func (u *CategoryUsecase) DeleteCategory(id string) error {
	return u.categoryRepo.DeleteCategory(id)
}
