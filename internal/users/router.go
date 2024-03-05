package users

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(
	r fiber.Router,
	tokenExp time.Duration,
	jwtService JWTService,
	userRepository UserRepository,
) {

	us := NewService(jwtService, userRepository)
	handle := newHandler(tokenExp, us)

	r.Post("register", handle.register)
	r.Post("login", handle.authenticate)
}
