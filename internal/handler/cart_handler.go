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

type CartItemRequest struct {
	ProductID string `json:"product_id" example:"prod-12345"`
	Quantity  int    `json:"quantity" example:"2"`
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

// GetCart godoc
// @Summary      Get user cart
// @Description  Retrieve the current user's shopping cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /user/cart [get]
// @Security     BearerAuth
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

// AddItem godoc
// @Summary      Add item to cart
// @Description  Add a product to the user's shopping cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request body CartItemRequest true "Cart Item (product_id, quantity)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /user/cart/items [post]
// @Security     BearerAuth
func (h *CartHandler) AddItem(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CartItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.cartUsecase.AddItem(userID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Item added to cart"})
}

// UpdateItem godoc
// @Summary      Update cart item
// @Description  Update the quantity of a product in the user's shopping cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request body CartItemRequest true "Cart Item Update (product_id, quantity)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /user/cart/items [put]
// @Security     BearerAuth
func (h *CartHandler) UpdateItem(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CartItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.cartUsecase.UpdateItem(userID, req.ProductID, req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Cart item updated"})
}

// RemoveItem godoc
// @Summary      Remove item from cart
// @Description  Remove a product from the user's shopping cart by product ID
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        productId path string true "Product ID"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /user/cart/items/{productId} [delete]
// @Security     BearerAuth
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

// ClearCart godoc
// @Summary      Clear cart
// @Description  Remove all items from the user's shopping cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /user/cart [delete]
// @Security     BearerAuth
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
