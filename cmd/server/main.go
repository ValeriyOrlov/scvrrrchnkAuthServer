package main

import (
	"strings"

	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/config"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/database"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/handler"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/repository"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverware "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type logrusWriter struct {
	logger *logrus.Logger
}

func (w *logrusWriter) Write(p []byte) (n int, err error) {
	message := strings.TrimSpace(string(p))
	w.logger.Info(message)
	return len(p), nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatalf("cannot load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.DatabaseDSN)
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to database")
	}

	if err := database.RunMigrations(db); err != nil {
		logrus.WithError(err).Fatal("migration failed")
	}

	userRepo := repository.NewGormUserRepo(db)
	tokenRepo := repository.NewGormTokenRepo(db)
	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
		cfg.JWTSecret,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(recoverware.New(recoverware.Config{
		EnableStackTrace: true,
	}))

	app.Use(logger.New(logger.Config{
		Output: &logrusWriter{logger: log},
	}))
	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.Refresh)

	app.Get("/me", handler.AuthRequired(cfg.JWTSecret), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"user_id": userID})
	})

	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		logrus.WithError(err).Fatal("server stopped")
	}
}
