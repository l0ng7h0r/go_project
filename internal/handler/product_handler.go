package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	sellerUsecase  *usecase.SellerUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase, sellerUsecase *usecase.SellerUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase, sellerUsecase: sellerUsecase}
}

func (h *ProductHandler) CreateProduct(c fiber.Ctx) error {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Stock       int      `json:"stock"`
		Status      string   `json:"status"`
		ImageURLs   []string `json:"image_urls"`
		CategoryIDs []string `json:"category_ids"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sellerID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	productID, err := h.productUsecase.CreateProduct(sellerID, req.Name, req.Description, req.Price, req.Stock, req.Status, req.ImageURLs, req.CategoryIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"product_id": productID, "message": "Product created successfully"})
}

func (h *ProductHandler) GetProductByID(c fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

func (h *ProductHandler) GetAllProducts(c fiber.Ctx) error {
	products, err := h.productUsecase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetProductsBySeller(c fiber.Ctx) error {
	id := c.Params("id")
	products, err := h.productUsecase.GetProductsBySeller(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")

	sellerID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.productUsecase.DeleteProduct(id, sellerID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

func (h *ProductHandler) UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Status      string  `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sellerID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.productUsecase.UpdateProduct(id, &domain.Product{
		SellerID:    sellerID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      req.Status,
	}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func (h *ProductHandler) GetProductsByCategory(c fiber.Ctx) error {
	categoryID := c.Params("id")
	products, err := h.productUsecase.GetProductsByCategory(categoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if products == nil {
		products = []domain.Product{}
	}
	return c.JSON(products)
}
