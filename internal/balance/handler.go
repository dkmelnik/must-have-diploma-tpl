package balance

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/dkmelnik/go-musthave-diploma/internal/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/logger"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
)

type (
	balanceService interface {
		GetCurrentBalance(ctx context.Context, userID models.ModelID) (dto.Balance, error)
	}
	handler struct {
		service balanceService
	}
)

func newHandler(ws balanceService) *handler {
	return &handler{ws}
}

func (h *handler) getCurrentBalance(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	out, err := h.service.GetCurrentBalance(c.Context(), models.ModelID(userID))
	if err != nil {
		logger.Log.Info("getCurrentBalance", "GetCurrentBalance", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(out)
}
