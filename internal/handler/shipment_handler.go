package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/golang/internal/usecase"
)

type ShipmentHandler struct {
	shipmentUsecase *usecase.ShipmentUsecase
}

func NewShipmentHandler(shipmentUsecase *usecase.ShipmentUsecase) *ShipmentHandler {
	return &ShipmentHandler{shipmentUsecase: shipmentUsecase}
}

func (h *ShipmentHandler) CreateShipment(c fiber.Ctx) error {
	var req struct {
		OrderID        string `json:"order_id"`
		Provider       string `json:"provider"`
		TrackingNumber string `json:"tracking_number"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	id, err := h.shipmentUsecase.CreateShipment(req.OrderID, req.Provider, req.TrackingNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"shipment_id": id, "message": "Shipment created"})
}

func (h *ShipmentHandler) GetShipmentByOrder(c fiber.Ctx) error {
	orderID := c.Params("orderId")
	shipment, err := h.shipmentUsecase.GetShipmentByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Shipment not found"})
	}
	return c.JSON(shipment)
}

func (h *ShipmentHandler) UpdateStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.shipmentUsecase.UpdateStatus(id, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Shipment status updated"})
}

func (h *ShipmentHandler) UpdateTracking(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Provider       string `json:"provider"`
		TrackingNumber string `json:"tracking_number"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.shipmentUsecase.UpdateTracking(id, req.Provider, req.TrackingNumber); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Tracking updated"})
}
