package main

import (
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/config"
	database "github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatalf("cannot load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.DatabaseDSN)
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to database")
	}
	_ = db
	app := fiber.New()

	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		logrus.WithError(err).Fatal("server stopped")
	}
}
