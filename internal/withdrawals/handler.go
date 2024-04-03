package withdrawals

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
	"github.com/dkmelnik/go-musthave-diploma/internal/utils"
	"github.com/dkmelnik/go-musthave-diploma/internal/withdrawals/dto"
)

type (
	withdrawalService interface {
		WithdrawAccrual(ctx context.Context, userID models.ModelID, d dto.WithdrawalPayload) error
		GetAllWithdrawals(ctx context.Context, userID models.ModelID) ([]dto.WithdrawalResponse, error)
	}
	handler struct {
		service withdrawalService
	}
)

func newHandler(ws withdrawalService) *handler {
	return &handler{ws}
}

func (h *handler) withdrawAccrual(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var body dto.WithdrawalPayload

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(http.StatusText(fiber.StatusUnprocessableEntity))
	}

	i, err := strconv.Atoi(body.Order)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if !utils.CheckNumberOnLuhn(i) {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	switch err = h.service.WithdrawAccrual(c.Context(), models.ModelID(userID), body); {
	case errors.Is(err, apperrors.ErrInsufficientFunds):
		return c.SendStatus(fiber.StatusPaymentRequired)
	case err == nil:
		return c.SendStatus(fiber.StatusOK)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func (h *handler) getAllWithdrawals(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	switch out, err := h.service.GetAllWithdrawals(c.Context(), models.ModelID(userID)); {
	case errors.Is(err, apperrors.ErrNoInformationAnswer):
		return c.SendStatus(fiber.StatusNoContent)
	case err == nil:
		return c.Status(fiber.StatusOK).JSON(out)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
