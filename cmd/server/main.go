package main

import (
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/config"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/database"
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

	if err := database.RunMigrations(db); err != nil {
		logrus.WithError(err).Fatal("migration failed")
	}

	app := fiber.New()

	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		logrus.WithError(err).Fatal("server stopped")
	}
}
