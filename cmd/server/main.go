package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"

    "furniture-shop/internal/config"
    "furniture-shop/internal/database"
    httpHandlers "furniture-shop/internal/adapter/http/handlers"
    httpRoutes "furniture-shop/internal/adapter/http/routes"
    infra "furniture-shop/internal/adapter/postgres"
    app "furniture-shop/internal/app"
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

    srv := fiber.New()
    srv.Use(cors.New(cors.Config{
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
    optRepo := infra.NewProductOptionRepository(db)
    orderRepo := infra.NewOrderRepository(db)

    authSvc := app.NewAuthService(userRepo, cfg.JWTSecret)
    catalogSvc := app.NewCatalogService(deptRepo, catRepo, prodRepo)
    ordersSvc := app.NewOrdersService(userRepo, orderRepo, prodRepo)
    adminSvc := app.NewAdminService(deptRepo, catRepo, prodRepo, optRepo)
    paymentsSvc := app.NewPaymentService(orderRepo)

    authHandlers := httpHandlers.NewAuthHandler(authSvc)
    catalogHandlers := httpHandlers.NewCatalogHandler(catalogSvc)
    ordersHandlers := httpHandlers.NewOrdersHandler(ordersSvc)
    adminHandlers := httpHandlers.NewAdminHandler(adminSvc)
    paymentsHandlers := httpHandlers.NewPaymentsHandler(paymentsSvc)

    httpRoutes.Register(srv, cfg, authHandlers, catalogHandlers, ordersHandlers, adminHandlers, paymentsHandlers)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server listening on :%s", port)
    log.Fatal(srv.Listen(":" + port))
}



