package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type CategoryHandler struct {
	categoryUsecase *usecase.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase}
}

type CategoryRequest struct {
	Name string `json:"name" example:"Electronics"`
}

// CreateCategory godoc
// @Summary      Create a new category (Admin)
// @Description  Create a new product category
// @Tags         admin, categories
// @Accept       json
// @Produce      json
// @Param        request body CategoryRequest true "Category details (name)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/categories/create [post]
// @Security     BearerAuth
func (h *CategoryHandler) CreateCategory(c fiber.Ctx) error {
	var req CategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	id, err := h.categoryUsecase.CreateCategory(req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id, "message": "Category created successfully"})
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Retrieve a list of all product categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /categories [get]
func (h *CategoryHandler) GetAllCategories(c fiber.Ctx) error {
	categories, err := h.categoryUsecase.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(categories)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Retrieve a product category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c fiber.Ctx) error {
	id := c.Params("id")
	category, err := h.categoryUsecase.GetCategoryByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(category)
}

// UpdateCategory godoc
// @Summary      Update a category (Admin)
// @Description  Update the name of an existing category
// @Tags         admin, categories
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID"
// @Param        request body CategoryRequest true "Category details (name)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/categories/update/{id} [put]
// @Security     BearerAuth
func (h *CategoryHandler) UpdateCategory(c fiber.Ctx) error {
	id := c.Params("id")
	var req CategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.categoryUsecase.UpdateCategory(id, req.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Category updated successfully"})
}

// DeleteCategory godoc
// @Summary      Delete a category (Admin)
// @Description  Delete a product category by its ID
// @Tags         admin, categories
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/categories/delete/{id} [delete]
// @Security     BearerAuth
func (h *CategoryHandler) DeleteCategory(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.categoryUsecase.DeleteCategory(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Category deleted successfully"})
}
