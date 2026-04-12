package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type AddressHandler struct {
	addressUsecase *usecase.AddressUsecase
}

func NewAddressHandler(addressUsecase *usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{addressUsecase: addressUsecase}
}

func (h *AddressHandler) CreateAddress(c fiber.Ctx) error {
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	var req struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Address    string `json:"address"`
		City       string `json:"city"`
		Country    string `json:"country"`
		PostalCode string `json:"postal_code"`
		IsDefault  bool   `json:"is_default"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := h.addressUsecase.CreateAddress(userID, req.Name, req.Phone, req.Address, req.City, req.Country, req.PostalCode, req.IsDefault)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id, "message": "Address created successfully"})
}

func (h *AddressHandler) GetMyAddresses(c fiber.Ctx) error {
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	addresses, err := h.addressUsecase.GetMyAddresses(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(addresses)
}

func (h *AddressHandler) UpdateAddress(c fiber.Ctx) error {
	id := c.Params("id")
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	var req domain.Address
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.UserID = userID

	if err := h.addressUsecase.UpdateAddress(id, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Address updated successfully"})
}

func (h *AddressHandler) DeleteAddress(c fiber.Ctx) error {
	id := c.Params("id")
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	if err := h.addressUsecase.DeleteAddress(id, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Address deleted successfully"})
}

func (h *AddressHandler) SetDefaultAddress(c fiber.Ctx) error {
	id := c.Params("id")
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, ok := userIDLocals.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user_id in token"})
	}

	if err := h.addressUsecase.SetDefaultAddress(id, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Default address set successfully"})
}
