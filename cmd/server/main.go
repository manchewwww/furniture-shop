package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"

    "furniture-shop/internal/config"
    "furniture-shop/internal/database"
    "furniture-shop/internal/handlers"
    infra "furniture-shop/internal/infrastructure/repository/gorm"
    "furniture-shop/internal/routes"
    "furniture-shop/internal/services"
)

func main() {
    _ = godotenv.Load()

    cfg := config.Load()
    if err := database.Connect(cfg); err != nil {
        log.Fatalf("DB connection failed: %v", err)
    }
    if err := database.AutoMigrateAndSeed(); err != nil {
        log.Fatalf("Migration/Seed failed: %v", err)
    }

    app := fiber.New()
    app.Use(cors.New(cors.Config{
        AllowOrigins:     cfg.CORSOrigins,
        AllowCredentials: true,
        AllowMethods:     "GET,POST,PATCH,DELETE,PUT",
        AllowHeaders:     "Authorization,Content-Type",
    }))

    // Wire dependencies (repositories -> services -> handlers)
    db := database.DB
    userRepo := infra.NewUserRepository(db)
    deptRepo := infra.NewDepartmentRepository(db)
    catRepo := infra.NewCategoryRepository(db)
    prodRepo := infra.NewProductRepository(db)

    authSvc := services.NewAuthService(userRepo, cfg.JWTSecret)
    catalogSvc := services.NewCatalogService(deptRepo, catRepo, prodRepo)

    authHandlers := handlers.NewAuthHandler(authSvc)
    catalogHandlers := handlers.NewCatalogHandler(catalogSvc)

    routes.Register(app, cfg, authHandlers, catalogHandlers)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server listening on :%s", port)
    log.Fatal(app.Listen(":" + port))
}

