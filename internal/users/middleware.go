package users

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
)

type (
	MiddlewareManager struct {
		jwtService JWTService
	}
)

func NewMiddlewareManager(jwtService JWTService) *MiddlewareManager {
	return &MiddlewareManager{jwtService}
}

func (m *MiddlewareManager) Auth(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).SendString(http.StatusText(fiber.StatusUnauthorized))
	}

	token, err := m.jwtService.ParseToken(cookie)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			return c.Status(fiber.StatusUnauthorized).SendString(http.StatusText(fiber.StatusUnauthorized))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(http.StatusText(fiber.StatusInternalServerError))
	}

	c.Locals("user_id", token.SUB)
	c.Locals("token", token)

	return c.Next()
}
