package config

import "os"

type Config struct {
    DBUser      string
    DBPass      string
    DBHost      string
    DBPort      string
    DBName      string
    CORSOrigins string
    JWTSecret   string
}

func getenv(k, def string) string {
    v := os.Getenv(k)
    if v == "" {
        return def
    }
    return v
}

func Load() *Config {
    return &Config{
        DBUser:      getenv("DB_USER", "postgres"),
        DBPass:      getenv("DB_PASS", "postgres"),
        DBHost:      getenv("DB_HOST", "127.0.0.1"),
        DBPort:      getenv("DB_PORT", "5432"),
        DBName:      getenv("DB_NAME", "furniture_shop"),
        CORSOrigins: getenv("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000"),
        JWTSecret:   getenv("JWT_SECRET", "dev-secret-change-me"),
    }
}

