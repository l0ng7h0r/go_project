package usecase

import (
	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type ShipmentUsecase struct {
	shipmentRepo *repository.ShipmentRepository
}

func NewShipmentUsecase(shipmentRepo *repository.ShipmentRepository) *ShipmentUsecase {
	return &ShipmentUsecase{shipmentRepo: shipmentRepo}
}

func (u *ShipmentUsecase) CreateShipment(orderID, provider, trackingNumber string) (string, error) {
	return u.shipmentRepo.CreateShipment(&domain.Shipment{
		OrderID:        orderID,
		Provider:       provider,
		TrackingNumber: trackingNumber,
		Status:         "pending",
	})
}

func (u *ShipmentUsecase) GetShipmentByOrderID(orderID string) (*domain.Shipment, error) {
	return u.shipmentRepo.GetShipmentByOrderID(orderID)
}

func (u *ShipmentUsecase) UpdateStatus(id, status string) error {
	return u.shipmentRepo.UpdateShipmentStatus(id, status)
}

func (u *ShipmentUsecase) UpdateTracking(id, provider, trackingNumber string) error {
	return u.shipmentRepo.UpdateTrackingNumber(id, provider, trackingNumber)
}
