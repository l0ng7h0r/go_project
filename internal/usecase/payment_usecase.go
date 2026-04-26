package usecase

import (
	"fmt"

	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
	"github.com/l0ng7h0r/golang/pkg/phajay"
)

type PaymentUsecase struct {
	paymentRepo  *repository.PaymentRepository
	orderRepo    *repository.OrderRepository
	phajayClient *phajay.Client
}

func NewPaymentUsecase(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository, phajayClient *phajay.Client) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo:  paymentRepo,
		orderRepo:    orderRepo,
		phajayClient: phajayClient,
	}
}

// CreatePayment creates a Phajay payment link and saves it to the DB
func (u *PaymentUsecase) CreatePayment(orderID string, amount float64) (string, string, error) {
	// Verify order exists
	order, err := u.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return "", "", fmt.Errorf("order not found: %w", err)
	}

	description := fmt.Sprintf("Order %s - %s", order.ID[:8], order.ReceiverName)

	// Call Phajay API to create payment link
	phajayResp, err := u.phajayClient.CreatePaymentLink(amount, description, orderID)
	if err != nil {
		return "", "", fmt.Errorf("failed to create phajay payment: %w", err)
	}

	payment := &domain.Payment{
		OrderID:    orderID,
		Method:     "phajay",
		Status:     "pending",
		Amount:     amount,
		PaymentURL: phajayResp.PaymentURL,
	}

	paymentID, err := u.paymentRepo.CreatePayment(payment)
	if err != nil {
		return "", "", err
	}

	return paymentID, phajayResp.PaymentURL, nil
}

func (u *PaymentUsecase) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	return u.paymentRepo.GetPaymentByOrderID(orderID)
}

// HandleWebhook processes a Phajay webhook callback
func (u *PaymentUsecase) HandleWebhook(payload *phajay.WebhookPayload) error {
	payment, err := u.paymentRepo.GetPaymentByOrderID(payload.OrderNo)
	if err != nil {
		return fmt.Errorf("payment not found for order: %s", payload.OrderNo)
	}

	switch payload.Status {
	case "success":
		if err := u.paymentRepo.UpdatePaymentStatus(payment.ID, "paid", payload.TransactionID); err != nil {
			return err
		}
		// Update order status to confirmed
		return u.orderRepo.UpdateOrderStatus(payload.OrderNo, "confirmed")
	case "failed", "cancelled":
		return u.paymentRepo.UpdatePaymentStatus(payment.ID, "failed", payload.TransactionID)
	default:
		return fmt.Errorf("unknown webhook status: %s", payload.Status)
	}
}

func (u *PaymentUsecase) ConfirmPayment(paymentID, transactionID string) error {
	return u.paymentRepo.UpdatePaymentStatus(paymentID, "paid", transactionID)
}
