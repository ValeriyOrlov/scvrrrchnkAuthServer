package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseDSN string
}

func Load() (*Config, error) {
	godotenv.Load()
	cfg := &Config{
		Port:        getEnv("APP_PORT", fmt.Sprintf(":%s", os.Getenv("APP_PORT"))),
		DatabaseDSN: getEnv("DB_DSN", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
