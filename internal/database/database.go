package database

import (
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
