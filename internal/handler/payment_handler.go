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

type CreatePaymentRequest struct {
	OrderID string  `json:"order_id" example:"ord-123456"`
	Amount  float64 `json:"amount" example:"1500.50"`
}

type ConfirmPaymentRequest struct {
	TransactionID string `json:"transaction_id" example:"txn-987654"`
}

// CreatePayment godoc
// @Summary      Create a payment
// @Description  Create a Phajay payment link for an order
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        request body CreatePaymentRequest true "Payment details (order_id, amount)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /user/payments [post]
// @Security     BearerAuth
func (h *PaymentHandler) CreatePayment(c fiber.Ctx) error {
	var req CreatePaymentRequest
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

// GetPaymentByOrder godoc
// @Summary      Get payment by order
// @Description  Get payment info for a specific order
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        orderId path string true "Order ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /user/payments/order/{orderId} [get]
// @Security     BearerAuth
func (h *PaymentHandler) GetPaymentByOrder(c fiber.Ctx) error {
	orderID := c.Params("orderId")
	payment, err := h.paymentUsecase.GetPaymentByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}
	return c.JSON(payment)
}

// PhajayWebhook godoc
// @Summary      Phajay payment webhook
// @Description  Handle the Phajay webhook callback (Public)
// @Tags         payments, webhooks
// @Accept       json
// @Produce      json
// @Param        request body phajay.WebhookPayload true "Webhook payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /webhooks/phajay [post]
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

// ConfirmPayment godoc
// @Summary      Confirm payment manually (Admin)
// @Description  Allows admin to manually confirm a payment
// @Tags         admin, payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID"
// @Param        request body ConfirmPaymentRequest true "Transaction details (transaction_id)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/payments/{id}/confirm [patch]
// @Security     BearerAuth
func (h *PaymentHandler) ConfirmPayment(c fiber.Ctx) error {
	id := c.Params("id")
	var req ConfirmPaymentRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.paymentUsecase.ConfirmPayment(id, req.TransactionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Payment confirmed"})
}
