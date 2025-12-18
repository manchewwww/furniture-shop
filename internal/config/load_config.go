package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Configurations Config

func LoadConfig() error {
	candidates := []string{}
	if p := strings.TrimSpace(os.Getenv("APP_CONFIG")); p != "" {
		candidates = append(candidates, p)
	}
	candidates = append(candidates,
		"appconfig.json",
		filepath.Join("internal", "config", "appconfig.json"),
		"/app/appconfig.json",
	)

	var file *os.File
	var tried []string
	for _, path := range candidates {
		f, err := os.Open(path)
		if err != nil {
			tried = append(tried, path)
			continue
		}
		file = f
		break
	}
	if file == nil {
		return fmt.Errorf("failed to open config file; tried: %s", strings.Join(tried, ", "))
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
