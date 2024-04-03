package users

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/users/dto"

	"github.com/gofiber/fiber/v2"
)

type (
	userService interface {
		Register(ctx context.Context, dto dto.RegisterPayload) (string, error)
		Authenticate(ctx context.Context, dto dto.LoginPayload) (string, error)
	}
	handler struct {
		tokenExp time.Duration
		service  userService
	}
)

func newHandler(tokenExp time.Duration, service userService) *handler {
	return &handler{tokenExp, service}
}

func (h *handler) setCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
		//HTTPOnly: true,
		Expires: time.Now().Add(h.tokenExp),
		//Secure:   true,
	})
}

func (h *handler) register(c *fiber.Ctx) error {
	var body dto.RegisterPayload

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(http.StatusText(fiber.StatusUnprocessableEntity))
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(http.StatusText(fiber.StatusUnprocessableEntity))
	}

	token, err := h.service.Register(c.Context(), body)
	if err != nil {
		log.Println(err)
		if errors.Is(err, apperrors.ErrIsExist) {
			return c.Status(fiber.StatusConflict).SendString(http.StatusText(fiber.StatusConflict))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(http.StatusText(fiber.StatusInternalServerError))
	}

	h.setCookie(c, token)

	return c.SendStatus(fiber.StatusOK)

}

func (h *handler) authenticate(c *fiber.Ctx) error {
	var body dto.LoginPayload

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(http.StatusText(fiber.StatusUnprocessableEntity))
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(http.StatusText(fiber.StatusUnprocessableEntity))
	}

	token, err := h.service.Authenticate(c.Context(), body)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return c.Status(fiber.StatusUnauthorized).SendString(http.StatusText(fiber.StatusUnauthorized))
		}
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).SendString(http.StatusText(fiber.StatusUnauthorized))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	h.setCookie(c, token)

	return c.SendStatus(fiber.StatusOK)
}
