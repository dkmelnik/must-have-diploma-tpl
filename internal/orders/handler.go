package orders

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/dkmelnik/go-musthave-diploma/internal/logger"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
	"github.com/dkmelnik/go-musthave-diploma/internal/orders/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/utils"
)

type (
	orderService interface {
		CreateIfNotExist(ctx context.Context, userID, orderNumber string) error
		GetAllUserOrders(ctx context.Context, userID models.ModelID) ([]dto.OrderResponse, error)
	}
	handler struct {
		service orderService
	}
)

func newHandler(service orderService) *handler {
	return &handler{service}
}

func (h *handler) create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	orderNumber := string(c.Body())

	if orderNumber == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	i, err := strconv.Atoi(orderNumber)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if !utils.CheckNumberOnLuhn(i) {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if err := h.service.CreateIfNotExist(c.Context(), userID, orderNumber); err != nil {
		switch {
		case strings.Contains(err.Error(), "order already exists for the user"):
			return c.SendStatus(fiber.StatusOK)
		case strings.Contains(err.Error(), "order already exists"):
			return c.SendStatus(fiber.StatusConflict)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func (h *handler) getAllOrders(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	orders, err := h.service.GetAllUserOrders(c.Context(), models.ModelID(userID))
	if err != nil {
		logger.Log.Error("orders:handler:getAllOrders", "StatusInternalServerError", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if len(orders) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Status(fiber.StatusOK).JSON(orders, "application/json")
}
