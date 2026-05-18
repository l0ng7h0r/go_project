package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

type CreateOrderRequest struct {
	ReceiverName   string `json:"receiver_name" example:"John Doe"`
	Phone          string `json:"phone" example:"0812345678"`
	Province       string `json:"province" example:"Bangkok"`
	District       string `json:"district" example:"Chatuchak"`
	Logistic       string `json:"logistic" example:"Kerry Express"`
	LogisticBranch string `json:"logistic_branch" example:"Branch 001"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"shipped"`
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Create an order from the user's current shopping cart
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body CreateOrderRequest true "Order details (receiver_name, phone, province, district, logistic, logistic_branch)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /user/orders [post]
// @Security     BearerAuth
func (h *OrderHandler) CreateOrder(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CreateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.ReceiverName == "" || req.Phone == "" || req.Province == "" || req.District == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "receiver_name, phone, province and district are required",
		})
	}

	orderID, err := h.orderUsecase.CreateOrderFromCart(
		userID,
		req.ReceiverName, req.Phone,
		req.Province, req.District,
		req.Logistic, req.LogisticBranch,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": orderID,
		"message":  "Order created successfully",
	})
}

// GetMyOrders godoc
// @Summary      Get user's orders
// @Description  Retrieve all orders placed by the current user
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /user/orders [get]
// @Security     BearerAuth
func (h *OrderHandler) GetMyOrders(c fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	orders, err := h.orderUsecase.GetMyOrders(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

// GetOrderByID godoc
// @Summary      Get order by ID
// @Description  Retrieve a specific order by its ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /user/orders/{id} [get]
// @Security     BearerAuth
func (h *OrderHandler) GetOrderByID(c fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.orderUsecase.GetOrderByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}
	return c.JSON(order)
}

// GetAllOrders godoc
// @Summary      Get all orders (Admin)
// @Description  Retrieve all orders in the system
// @Tags         admin, orders
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/orders [get]
// @Security     BearerAuth
func (h *OrderHandler) GetAllOrders(c fiber.Ctx) error {
	orders, err := h.orderUsecase.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

// UpdateOrderStatus godoc
// @Summary      Update order status (Admin)
// @Description  Update the status of an existing order
// @Tags         admin, orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Param        request body UpdateOrderStatusRequest true "Status Update (status)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /admin/orders/{id}/status [patch]
// @Security     BearerAuth
func (h *OrderHandler) UpdateOrderStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateOrderStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.orderUsecase.UpdateOrderStatus(id, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order status updated"})
}
