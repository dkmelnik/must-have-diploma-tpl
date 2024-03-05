package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app  *fiber.App
	addr string
}

func NewServer(addr string) *Server {
	return &Server{fiber.New(), addr}
}

func (s *Server) Run() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	go func() {
		<-sig
		err := s.app.Shutdown()
		if err != nil {
			panic(err)
		}
	}()

	return s.app.Listen(s.addr)
}

func (s *Server) GetApp() *fiber.App {
	return s.app
}
