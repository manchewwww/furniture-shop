package routes

import (
    "github.com/gofiber/fiber/v2"

    "furniture-shop/internal/config"
    "furniture-shop/internal/handlers"
    "furniture-shop/internal/middleware"
)

func Register(app *fiber.App, cfg *config.Config, auth *handlers.AuthHandler, catalog *handlers.CatalogHandler, orders *handlers.OrdersHandler, admin *handlers.AdminHandler, payments *handlers.PaymentsHandler) {
    api := app.Group("/api")

    // Auth
    api.Post("/auth/register", auth.Register())
    api.Post("/auth/login", auth.Login())

    // Catalog & search
    api.Get("/departments", catalog.GetDepartments())
    api.Get("/departments/:id/categories", catalog.GetCategoriesByDepartment())
    api.Get("/categories/:id/products", catalog.GetProductsByCategory())
    api.Get("/products/:id", catalog.GetProductDetails())
    api.Get("/products/:id/recommendations", catalog.GetProductRecommendations())
    api.Get("/products/search", catalog.SearchProducts())

    // Orders (public create)
    api.Post("/orders", orders.CreateOrder())
    api.Post("/payments/card", payments.PayByCard())

    // Authenticated user routes
    authGroup := api.Group("/user", middleware.JWTAuth(cfg))
    authGroup.Get("/me", auth.Me())
    authGroup.Get("/orders", orders.UserOrders())
    authGroup.Get("/orders/:id", orders.UserOrderDetails())

    // Admin routes
    adminGroup := api.Group("/admin", middleware.JWTAuth(cfg), middleware.RequireAdmin)
    adminGroup.Get("/departments", admin.ListDepartments())
    adminGroup.Post("/departments", admin.CreateDepartment())
    adminGroup.Put("/departments/:id", admin.UpdateDepartment())
    adminGroup.Delete("/departments/:id", admin.DeleteDepartment())

    adminGroup.Get("/categories", admin.ListCategories())
    adminGroup.Post("/categories", admin.CreateCategory())
    adminGroup.Put("/categories/:id", admin.UpdateCategory())
    adminGroup.Delete("/categories/:id", admin.DeleteCategory())

    adminGroup.Get("/products", admin.ListProducts())
    adminGroup.Post("/products", admin.CreateProduct())
    adminGroup.Put("/products/:id", admin.UpdateProduct())
    adminGroup.Delete("/products/:id", admin.DeleteProduct())

    adminGroup.Get("/product_options", admin.ListProductOptions())
    adminGroup.Post("/product_options", admin.CreateProductOption())
    adminGroup.Put("/product_options/:id", admin.UpdateProductOption())
    adminGroup.Delete("/product_options/:id", admin.DeleteProductOption())

    adminGroup.Get("/orders", orders.AdminListOrders())
    adminGroup.Patch("/orders/:id/status", orders.AdminUpdateOrderStatus())
}

