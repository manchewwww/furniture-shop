package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Env EnvConfig

func LoadEnvFile() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		return fmt.Errorf("STRIPE_SECRET_KEY is required")
	}

	Env = EnvConfig{
		DBUser:          dbUser,
		DBPass:          dbPass,
		JWTSecret:       jwtSecret,
		StripeSecretKey: stripeSecretKey,
	}

	return nil
}
