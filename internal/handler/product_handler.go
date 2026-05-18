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

type CreateProductRequest struct {
	Name        string   `json:"name" example:"Gaming Mouse"`
	Description string   `json:"description" example:"High precision RGB gaming mouse"`
	Price       float64  `json:"price" example:"1250.00"`
	Stock       int      `json:"stock" example:"50"`
	Status      string   `json:"status" example:"active"`
	ImageURLs   []string `json:"image_urls" example:"http://example.com/img1.png,http://example.com/img2.png"`
	CategoryIDs []string `json:"category_ids" example:"cat-1,cat-2"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" example:"Gaming Mouse V2"`
	Description string  `json:"description" example:"Upgraded version"`
	Price       float64 `json:"price" example:"1500.00"`
	Stock       int     `json:"stock" example:"30"`
	Status      string  `json:"status" example:"active"`
}

// CreateProduct godoc
// @Summary      Create a new product (Seller)
// @Description  Create a new product listed by the current seller
// @Tags         products, seller
// @Accept       json
// @Produce      json
// @Param        request body CreateProductRequest true "Product details"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /seller/products/create [post]
// @Security     BearerAuth
func (h *ProductHandler) CreateProduct(c fiber.Ctx) error {
	var req CreateProductRequest
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

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Retrieve details of a specific product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProductByID(c fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /products [get]
func (h *ProductHandler) GetAllProducts(c fiber.Ctx) error {
	products, err := h.productUsecase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

// GetProductsBySeller godoc
// @Summary      Get products by seller
// @Description  Retrieve all products listed by a specific seller
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Seller ID"
// @Success      200 {array} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /products/seller/{id} [get]
func (h *ProductHandler) GetProductsBySeller(c fiber.Ctx) error {
	id := c.Params("id")
	products, err := h.productUsecase.GetProductsBySeller(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

// DeleteProduct godoc
// @Summary      Delete a product (Seller)
// @Description  Delete a product listed by the current seller
// @Tags         products, seller
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /seller/products/delete/{id} [delete]
// @Security     BearerAuth
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

// UpdateProduct godoc
// @Summary      Update a product (Seller)
// @Description  Update details of a product listed by the current seller
// @Tags         products, seller
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID"
// @Param        request body UpdateProductRequest true "Updated product details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /seller/products/update/{id} [put]
// @Security     BearerAuth
func (h *ProductHandler) UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateProductRequest
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

// GetProductsByCategory godoc
// @Summary      Get products by category
// @Description  Retrieve all products in a specific category
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID"
// @Success      200 {array} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /products/category/{id} [get]
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
