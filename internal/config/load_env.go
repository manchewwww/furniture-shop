package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Env EnvConfig

func LoadEnvFile() error {
	err := godotenv.Load(".env")
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return nil
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	Env = EnvConfig{
		DBUser:    dbUser,
		DBPass:    dbPass,
		JWTSecret: jwtSecret,
	}

	return fmt.Errorf("failed to load .env file: %w", err)
}
