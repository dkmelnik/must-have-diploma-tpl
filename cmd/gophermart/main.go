package main

import (
	"database/sql"
	"errors"
	"fmt"

	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"github.com/dkmelnik/go-musthave-diploma/configs"
	"github.com/dkmelnik/go-musthave-diploma/internal/balance"
	"github.com/dkmelnik/go-musthave-diploma/internal/db/pg"
	"github.com/dkmelnik/go-musthave-diploma/internal/jwt"
	"github.com/dkmelnik/go-musthave-diploma/internal/logger"
	"github.com/dkmelnik/go-musthave-diploma/internal/orders"
	"github.com/dkmelnik/go-musthave-diploma/internal/server"
	"github.com/dkmelnik/go-musthave-diploma/internal/users"
	"github.com/dkmelnik/go-musthave-diploma/internal/withdrawals"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	godotenv.Load()

	conf, err := configs.NewServer()
	if err != nil {
		return err
	}

	// LOGGER -----------------------
	logger.Setup(conf.LogLevel, os.Stdout)
	// LOGGER -----------------------

	// PG -----------------------
	pgConnection, err := pg.NewConnection(conf.PGUri)
	if err != nil {
		logger.Log.Error("run", "pgConnection err:", err)
	}
	if err = migrateDB(conf.PGUri); err != nil {
		logger.Log.Error("run", "migrateDB err:", err)
	}
	// PG -----------------------

	// SERVER -----------------------
	srv := server.NewServer(conf.ServerAddr)

	if err = setupRouting(conf, srv.GetApp(), pgConnection); err != nil {
		return err
	}

	return srv.Run()
}

func setupRouting(conf configs.Server, s *fiber.App, db *sql.DB) error {
	api := s.Group("/api/user")
	api.Use(requestid.New())
	api.Use(fiberlogger.New())
	api.Use(recover.New())
	exp := time.Duration(conf.JWTExp) * time.Hour

	// repositories
	userRepository := users.NewRepository(db)
	orderRepository := orders.NewRepository(db)
	withdrawalRepository := withdrawals.NewRepository(db)

	//infrastructure services
	jwtService := jwt.NewJwt(conf.JWTSecret, exp)
	userMiddleware := users.NewMiddlewareManager(jwtService)
	balanceService := balance.NewService(withdrawalRepository, orderRepository)

	users.SetupRouter(api, exp, jwtService, userRepository)
	orders.SetupRouter(api, conf.AccrualAddr, userMiddleware, orderRepository)
	withdrawals.SetupRouter(api, userMiddleware, balanceService, withdrawalRepository)
	balance.SetupRouter(api, userMiddleware, balanceService)

	return nil
}

func migrateDB(dsn string) error {
	m, err := migrate.New("file://internal/db/pg/migrate", dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
