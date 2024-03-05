package balance

import (
	"github.com/gofiber/fiber/v2"
)

type userMiddleware interface {
	Auth(c *fiber.Ctx) error
}

func SetupRouter(
	r fiber.Router,
	middleware userMiddleware,
	bs balanceService,
) {
	handle := newHandler(bs)

	r.Get("balance", middleware.Auth, handle.getCurrentBalance)
}
