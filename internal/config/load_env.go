package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvFile() error {
	err := godotenv.Load(".env")
	if err == nil {
		return nil
	}

	// Ignore missing .env file (common in production where env vars come from the system).
	if os.IsNotExist(err) {
		return nil
	}

	return fmt.Errorf("failed to load .env file: %w", err)
}
