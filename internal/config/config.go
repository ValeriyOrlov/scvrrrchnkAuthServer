package config

import "os"

type Config struct {
	Port string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port: getEnv("APP_PORT", ":8080"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
