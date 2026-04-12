package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type CartHandler struct {
	cartUsecase *usecase.CartUsecase
}

func NewCartHandler(cartUsecase *usecase.CartUsecase) *CartHandler {
	return &CartHandler{cartUsecase: cartUsecase}
}

func getUserIDFromLocals(c fiber.Ctx) (string, error) {
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return "", fiber.ErrUnauthorized
	}
	userID, ok := userIDLocals.(string)
	if !ok {
		return "", fiber.ErrInternalServerError
	}
	return userID, nil
}

func (h *CartHandler) GetCart(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	cart, err := h.cartUsecase.GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cart)
}

func (h *CartHandler) AddItem(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.cartUsecase.AddItem(userID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Item added to cart"})
}

func (h *CartHandler) UpdateItem(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.cartUsecase.UpdateItem(userID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Cart item updated"})
}

func (h *CartHandler) RemoveItem(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	productID := c.Params("productId")
	if err := h.cartUsecase.RemoveItem(userID, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}

func (h *CartHandler) ClearCart(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if err := h.cartUsecase.ClearCart(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Cart cleared"})
}
