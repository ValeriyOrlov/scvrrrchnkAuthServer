package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Fatalf("cannot load db: %v", err)
	}
	return db, nil
}
