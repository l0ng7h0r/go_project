package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDsn            string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	AppPort          string
	PhajaySecretKey  string
}

var cfg *Config

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
	}
	return val
}

func Load() *Config{
	if cfg != nil {
		return cfg
	}

	_ = godotenv.Load()

	cfg = &Config{
		DBDsn:           getEnv("DB_DSN"),
		JWTSecret:       getEnv("JWT_SECRET"),
		JWTAccessExpiry: getEnv("JWT_ACCESS_EXPIRY"),
		JWTRefreshExpiry:getEnv("JWT_REFRESH_EXPIRY"),
		AppPort:         getEnv("APP_PORT"),
		PhajaySecretKey: os.Getenv("PHAJAY_SECRET_KEY"),
	}

	return cfg
}