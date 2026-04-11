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

func (h *SellerHandler) CreateSeller(c fiber.Ctx) error {

	var req struct {
		Email       string   `json:"email"`
		Password    string   `json:"password"`
		Roles       []string `json:"roles"`
		StoreName   string   `json:"store_name"`
		Description string   `json:"description"`
	}

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

func (h *SellerHandler) GetSellerByID(c fiber.Ctx) error {
	id := c.Params("id")
	seller, err := h.sellerUsecase.GetSellerByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(seller)
}

func (h *SellerHandler) GetAllSellers(c fiber.Ctx) error {
	sellers, err := h.sellerUsecase.GetAllSellers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sellers)
}

func (h *SellerHandler) DeleteSeller(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.sellerUsecase.DeleteSeller(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Seller deleted successfully"})
}

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