package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var Configurations Config

func LoadConfig() error {
	primaryPath := "/app/appconfig.json"

	file, err := os.Open(primaryPath)
	if err != nil {
		return fmt.Errorf("failed to read both primary and fallback config files: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read config file content: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	Configurations = cfg
	return nil
}
