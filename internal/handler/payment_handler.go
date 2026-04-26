package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
	"github.com/l0ng7h0r/golang/pkg/phajay"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}

// CreatePayment creates a Phajay payment link for an order
func (h *PaymentHandler) CreatePayment(c fiber.Ctx) error {
	var req struct {
		OrderID string  `json:"order_id"`
		Amount  float64 `json:"amount"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.OrderID == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "order_id and amount are required"})
	}

	paymentID, paymentURL, err := h.paymentUsecase.CreatePayment(req.OrderID, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"payment_id":  paymentID,
		"payment_url": paymentURL,
		"message":     "Payment created — redirect customer to payment_url",
	})
}

// GetPaymentByOrder gets payment info for an order
func (h *PaymentHandler) GetPaymentByOrder(c fiber.Ctx) error {
	orderID := c.Params("orderId")
	payment, err := h.paymentUsecase.GetPaymentByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}
	return c.JSON(payment)
}

// PhajayWebhook handles the Phajay webhook callback (Public — Phajay calls this)
func (h *PaymentHandler) PhajayWebhook(c fiber.Ctx) error {
	var payload phajay.WebhookPayload
	if err := c.Bind().Body(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.paymentUsecase.HandleWebhook(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "webhook processed"})
}

// ConfirmPayment allows admin to manually confirm a payment
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
