package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}

func (h *PaymentHandler) CreatePayment(c fiber.Ctx) error {
	var req struct {
		OrderID string  `json:"order_id"`
		Method  string  `json:"method"`
		Amount  float64 `json:"amount"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	id, err := h.paymentUsecase.CreatePayment(req.OrderID, req.Method, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"payment_id": id, "message": "Payment created"})
}

func (h *PaymentHandler) GetPaymentByOrder(c fiber.Ctx) error {
	orderID := c.Params("orderId")
	payment, err := h.paymentUsecase.GetPaymentByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}
	return c.JSON(payment)
}

func (h *PaymentHandler) ConfirmPayment(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		TransactionID string `json:"transaction_id"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.paymentUsecase.ConfirmPayment(id, req.TransactionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Payment confirmed"})
}
