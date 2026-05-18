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

type CreateShipmentRequest struct {
	OrderID        string `json:"order_id" example:"ord-123456"`
	Provider       string `json:"provider" example:"Kerry Express"`
	TrackingNumber string `json:"tracking_number" example:"KER123456789TH"`
}

type UpdateShipmentStatusRequest struct {
	Status string `json:"status" example:"in_transit"`
}

type UpdateShipmentTrackingRequest struct {
	Provider       string `json:"provider" example:"Kerry Express"`
	TrackingNumber string `json:"tracking_number" example:"KER987654321TH"`
}

// CreateShipment godoc
// @Summary      Create a shipment (Admin)
// @Description  Create a shipment for an order
// @Tags         admin, shipments
// @Accept       json
// @Produce      json
// @Param        request body CreateShipmentRequest true "Shipment details (order_id, provider, tracking_number)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/shipments/create [post]
// @Security     BearerAuth
func (h *ShipmentHandler) CreateShipment(c fiber.Ctx) error {
	var req CreateShipmentRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	id, err := h.shipmentUsecase.CreateShipment(req.OrderID, req.Provider, req.TrackingNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"shipment_id": id, "message": "Shipment created"})
}

// GetShipmentByOrder godoc
// @Summary      Get shipment by order
// @Description  Retrieve shipment information for a specific order
// @Tags         shipments
// @Accept       json
// @Produce      json
// @Param        orderId path string true "Order ID"
// @Success      200 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /user/shipments/order/{orderId} [get]
// @Security     BearerAuth
func (h *ShipmentHandler) GetShipmentByOrder(c fiber.Ctx) error {
	orderID := c.Params("orderId")
	shipment, err := h.shipmentUsecase.GetShipmentByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Shipment not found"})
	}
	return c.JSON(shipment)
}

// UpdateStatus godoc
// @Summary      Update shipment status (Admin)
// @Description  Update the status of a shipment
// @Tags         admin, shipments
// @Accept       json
// @Produce      json
// @Param        id path string true "Shipment ID"
// @Param        request body UpdateShipmentStatusRequest true "Status update (status)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/shipments/{id}/status [patch]
// @Security     BearerAuth
func (h *ShipmentHandler) UpdateStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateShipmentStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.shipmentUsecase.UpdateStatus(id, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Shipment status updated"})
}

// UpdateTracking godoc
// @Summary      Update shipment tracking (Admin)
// @Description  Update the tracking number and provider for a shipment
// @Tags         admin, shipments
// @Accept       json
// @Produce      json
// @Param        id path string true "Shipment ID"
// @Param        request body UpdateShipmentTrackingRequest true "Tracking info (provider, tracking_number)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /admin/shipments/{id}/tracking [patch]
// @Security     BearerAuth
func (h *ShipmentHandler) UpdateTracking(c fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateShipmentTrackingRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.shipmentUsecase.UpdateTracking(id, req.Provider, req.TrackingNumber); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Tracking updated"})
}
