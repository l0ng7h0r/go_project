package handler


import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type SellerHandler struct {
	sellerUsecase *usecase.SellerUsecase
}

func NewSellerHandler(sellerUsecase *usecase.SellerUsecase) *SellerHandler {
	return &SellerHandler{sellerUsecase: sellerUsecase}
}

type CreateSellerRequest struct {
	Email       string   `json:"email" example:"seller@example.com"`
	Password    string   `json:"password" example:"seller123"`
	Roles       []string `json:"roles" example:"seller"`
	StoreName   string   `json:"store_name" example:"Cool Gadgets Store"`
	Description string   `json:"description" example:"We sell the coolest gadgets around."`
}

// CreateSeller godoc
// @Summary      Create a new seller (Admin)
// @Description  Create a new seller account and store
// @Tags         admin, sellers
// @Accept       json
// @Produce      json
// @Param        request body CreateSellerRequest true "Seller details (email, password, roles, store_name, description)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/sellers [post]
// @Security     BearerAuth
func (h *SellerHandler) CreateSeller(c fiber.Ctx) error {

	var req CreateSellerRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	seller := &domain.Seller{
		StoreName:   req.StoreName,
		Description: req.Description,
	}

	if err := h.sellerUsecase.CreateSeller(req.Email, req.Password, req.Roles, seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Seller created successfully"})
}

// GetSellerByID godoc
// @Summary      Get seller by ID (Admin)
// @Description  Retrieve details of a specific seller
// @Tags         admin, sellers
// @Accept       json
// @Produce      json
// @Param        id path string true "Seller ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/sellers/{id} [get]
// @Security     BearerAuth
func (h *SellerHandler) GetSellerByID(c fiber.Ctx) error {
	id := c.Params("id")
	seller, err := h.sellerUsecase.GetSellerByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(seller)
}

// GetAllSellers godoc
// @Summary      Get all sellers (Admin)
// @Description  Retrieve a list of all sellers
// @Tags         admin, sellers
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/sellers [get]
// @Security     BearerAuth
func (h *SellerHandler) GetAllSellers(c fiber.Ctx) error {
	sellers, err := h.sellerUsecase.GetAllSellers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sellers)
}

// DeleteSeller godoc
// @Summary      Delete a seller (Admin)
// @Description  Delete a seller account
// @Tags         admin, sellers
// @Accept       json
// @Produce      json
// @Param        id path string true "Seller ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/sellers/{id} [delete]
// @Security     BearerAuth
func (h *SellerHandler) DeleteSeller(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.sellerUsecase.DeleteSeller(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Seller deleted successfully"})
}

// UpdateSeller godoc
// @Summary      Update a seller (Admin)
// @Description  Update details of an existing seller
// @Tags         admin, sellers
// @Accept       json
// @Produce      json
// @Param        id path string true "Seller ID"
// @Param        request body domain.Seller true "Updated seller details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/sellers/{id} [put]
// @Security     BearerAuth
func (h *SellerHandler) UpdateSeller(c fiber.Ctx) error {
	id := c.Params("id")
	var seller domain.Seller
	if err := c.Bind().Body(&seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.sellerUsecase.UpdateSeller(id, &seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Seller updated successfully"})
}