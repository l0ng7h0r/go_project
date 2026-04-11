package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	sellerUsecase *usecase.SellerUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase, sellerUsecase *usecase.SellerUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase, sellerUsecase: sellerUsecase}
}

func (h *ProductHandler) CreateProduct(c fiber.Ctx) error {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	sellerID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	err := h.productUsecase.CreateProduct(sellerID, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product created successfully"})
}

func (h *ProductHandler) GetProductByID(c fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

func (h *ProductHandler) GetAllProducts(c fiber.Ctx) error {
	products, err := h.productUsecase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetProductsBySeller(c fiber.Ctx) error {
	id := c.Params("id")
	products, err := h.productUsecase.GetProductsBySeller(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")

	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	sellerID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	err := h.productUsecase.DeleteProduct(id, sellerID)
	if err != nil {
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
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	sellerID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	err := h.productUsecase.UpdateProduct(id, &domain.Product{
		SellerID:    sellerID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

