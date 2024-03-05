package withdrawals

import (
	"github.com/gofiber/fiber/v2"
)

type UserMiddleware interface {
	Auth(c *fiber.Ctx) error
}

func SetupRouter(
	r fiber.Router,
	mw UserMiddleware,
	bs balanceService,
	wr withdrawalRepository,
) {
	us := NewService(bs, wr)
	handle := newHandler(us)

	r.Post("balance/withdraw", mw.Auth, handle.withdrawAccrual)
	r.Get("withdrawals", mw.Auth, handle.getAllWithdrawals)
}
