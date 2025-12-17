package config

import "github.com/joho/godotenv"

func LoadEnvFile() error {
	err := godotenv.Load(".env")
	return err
}
