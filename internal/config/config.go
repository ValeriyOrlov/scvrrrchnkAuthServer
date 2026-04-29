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

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var dbCfg DBConfig
	dbCfg.Host = getEnv("DB_HOST", "localhost")
	dbCfg.Port = getEnv("DB_PORT", "5432")
	dbCfg.User = getEnv("DB_USER", "postgres")
	dbCfg.Password = getEnv("DB_PASSWORD", "")
	dbCfg.DBName = getEnv("DB_NAME", "auth")
	dbCfg.SSLMode = getEnv("DB_SSLMODE", "disable")
	cfg := &Config{
		Port:        getEnv("APP_PORT", "8080"),
		DatabaseDSN: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Host, dbCfg.Port, dbCfg.SSLMode),
	}
	cfg.Port = ":" + cfg.Port
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
