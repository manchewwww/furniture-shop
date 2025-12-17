package main

import (
    "log"
    "os"

    "furniture-shop/internal/config"
    "furniture-shop/internal/database"
    httpserver "furniture-shop/internal/server/http"
    domain "furniture-shop/internal/service/domain"
    pg "furniture-shop/internal/storage/postgres"
)

func main() {
    // Load .env file first (if it exists) before loading config
    if err := config.LoadEnvFile(); err != nil { log.Fatalf("Env load failed: %v", err) }
    if err := config.LoadConfig(); err != nil { log.Fatalf("Config load failed: %v", err) }
    if err := database.Connect(); err != nil { log.Fatalf("DB connection failed: %v", err) }
    if err := database.AutoMigrateAndSeed(); err != nil { log.Fatalf("Migration/Seed failed: %v", err) }

    repos := pg.NewRepository(database.DB)
    jwtSecret := os.Getenv("JWT_SECRET")
    svc := domain.NewService(repos, jwtSecret)
    srv := httpserver.NewServer(svc)
    log.Fatal(srv.Run())
}

