package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type PaymentUsecase struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
}

func NewPaymentUsecase(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository) *PaymentUsecase {
	return &PaymentUsecase{paymentRepo: paymentRepo, orderRepo: orderRepo}
}

func (u *PaymentUsecase) CreatePayment(orderID, method string, amount float64) (string, error) {
	payment := &domain.Payment{
		OrderID: orderID,
		Method:  method,
		Status:  "pending",
		Amount:  amount,
	}
	return u.paymentRepo.CreatePayment(payment)
}

func (u *PaymentUsecase) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	return u.paymentRepo.GetPaymentByOrderID(orderID)
}

func (u *PaymentUsecase) ConfirmPayment(paymentID, transactionID string) error {
	return u.paymentRepo.UpdatePaymentStatus(paymentID, "paid", transactionID)
}

func (u *PaymentUsecase) FailPayment(paymentID string) error {
	return u.paymentRepo.UpdatePaymentStatus(paymentID, "failed", "")
}
