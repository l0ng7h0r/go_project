package usecase

import (
	"errors"

	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type ShipmentUsecase struct {
	shipmentRepo *repository.ShipmentRepository
	orderRepo    *repository.OrderRepository
}

func NewShipmentUsecase(shipmentRepo *repository.ShipmentRepository, orderRepo *repository.OrderRepository) *ShipmentUsecase {
	return &ShipmentUsecase{shipmentRepo: shipmentRepo, orderRepo: orderRepo}
}

func (u *ShipmentUsecase) CreateShipment(orderID, provider, trackingNumber string) (string, error) {
	return u.shipmentRepo.CreateShipment(&domain.Shipment{
		OrderID:        orderID,
		Provider:       provider,
		TrackingNumber: trackingNumber,
		Status:         "pending",
	})
}

func (u *ShipmentUsecase) GetShipmentByOrderID(orderID, currentUserID string, userRoles []interface{}) (*domain.Shipment, error) {
	// Verify the caller owns the order (or is an admin) before exposing shipment details
	order, err := u.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != currentUserID && !hasShipmentRole(userRoles, "admin") {
		return nil, errors.New("forbidden: access denied to this shipment")
	}
	return u.shipmentRepo.GetShipmentByOrderID(orderID)
}

func hasShipmentRole(roles []interface{}, target string) bool {
	for _, r := range roles {
		if s, ok := r.(string); ok && s == target {
			return true
		}
	}
	return false
}

func (u *ShipmentUsecase) UpdateStatus(id, status string) error {
	return u.shipmentRepo.UpdateShipmentStatus(id, status)
}

func (u *ShipmentUsecase) UpdateTracking(id, provider, trackingNumber string) error {
	return u.shipmentRepo.UpdateTrackingNumber(id, provider, trackingNumber)
}
