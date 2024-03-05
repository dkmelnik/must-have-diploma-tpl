package orders

import (
	"github.com/gofiber/fiber/v2"
)

type UserMiddleware interface {
	Auth(c *fiber.Ctx) error
}

func SetupRouter(
	r fiber.Router,
	accrualAddr string,
	middleware UserMiddleware,
	orderRepository orderRepository,
) {
	group := r.Group("/orders")

	wr := newWorker(accrualAddr, orderRepository)
	service := NewService(wr, orderRepository)
	handle := newHandler(service)

	group.Post("/", middleware.Auth, handle.create)
	group.Get("/", middleware.Auth, handle.getAllOrders)

}
