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
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if stripeWebhookSecret == "" {
		return fmt.Errorf("STRIPE_WEBHOOK_SECRET is required")
	}
	emailSenderHost := os.Getenv("SMTP_HOST")
	if emailSenderHost == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	emailSenderPort := os.Getenv("SMTP_PORT")
	if emailSenderPort == "" {
		return fmt.Errorf("SMTP_PORT is required")
	}
	emailSenderUser := os.Getenv("SMTP_USER")
	if emailSenderUser == "" {
		return fmt.Errorf("SMTP_USER is required")
	}
	emailSenderPass := os.Getenv("SMTP_PASS")
	if emailSenderPass == "" {
		return fmt.Errorf("SMTP_PASS is required")
	}
	emailSenderFrom := os.Getenv("FROM_EMAIL")
	if emailSenderFrom == "" {
		return fmt.Errorf("FROM_EMAIL is required")
	}

	Env = EnvConfig{
		DBUser:              dbUser,
		DBPass:              dbPass,
		JWTSecret:           jwtSecret,
		StripeSecretKey:     stripeSecretKey,
		StripeWebhookSecret: stripeWebhookSecret,
	}

	return nil
}
