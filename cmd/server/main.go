package main

import (
	"log"

	"furniture-shop/internal/config"
	"furniture-shop/internal/database"
	httpserver "furniture-shop/internal/server/http"
	domain "furniture-shop/internal/service/domain"
	pg "furniture-shop/internal/storage/postgres"
)

func main() {
	if err := config.LoadEnvFile(); err != nil {
		log.Fatalf("Env load failed: %v", err)
	}
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Config load failed: %v", err)
	}
	if err := database.Connect(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	if err := database.AutoMigrateAndSeed(); err != nil {
		log.Fatalf("Migration/Seed failed: %v", err)
	}

	repos := pg.NewRepository(database.DB)
	svc := domain.NewService(repos, config.Env.JWTSecret)
	srv := httpserver.NewServer(svc)
	log.Fatal(srv.Run())
}
