package routes

import (
    "github.com/gofiber/fiber/v2"

    "furniture-shop/internal/config"
    "furniture-shop/internal/handlers"
    "furniture-shop/internal/middleware"
)

func Register(app *fiber.App, cfg *config.Config) {
    api := app.Group("/api")

    // Auth
    api.Post("/auth/register", handlers.Register(cfg))
    api.Post("/auth/login", handlers.Login(cfg))

    // Catalog & search
    api.Get("/departments", handlers.GetDepartments)
    api.Get("/departments/:id/categories", handlers.GetCategoriesByDepartment)
    api.Get("/categories/:id/products", handlers.GetProductsByCategory)
    api.Get("/products/:id", handlers.GetProductDetails)
    api.Get("/products/:id/recommendations", handlers.GetProductRecommendations)
    api.Get("/products/search", handlers.SearchProducts)

    // Orders (public create)
    api.Post("/orders", handlers.CreateOrder)
    api.Post("/payments/card", handlers.PayByCard)

    // Authenticated user routes
    auth := api.Group("/user", middleware.JWTAuth(cfg))
    auth.Get("/me", handlers.Me())
    auth.Get("/orders", handlers.UserOrders)
    auth.Get("/orders/:id", handlers.UserOrderDetails)

    // Admin routes
    admin := api.Group("/admin", middleware.JWTAuth(cfg), middleware.RequireAdmin)
    admin.Get("/departments", handlers.AdminListDepartments)
    admin.Post("/departments", handlers.AdminCreateDepartment)
    admin.Put("/departments/:id", handlers.AdminUpdateDepartment)
    admin.Delete("/departments/:id", handlers.AdminDeleteDepartment)

    admin.Get("/categories", handlers.AdminListCategories)
    admin.Post("/categories", handlers.AdminCreateCategory)
    admin.Put("/categories/:id", handlers.AdminUpdateCategory)
    admin.Delete("/categories/:id", handlers.AdminDeleteCategory)

    admin.Get("/products", handlers.AdminListProducts)
    admin.Post("/products", handlers.AdminCreateProduct)
    admin.Put("/products/:id", handlers.AdminUpdateProduct)
    admin.Delete("/products/:id", handlers.AdminDeleteProduct)

    admin.Get("/product_options", handlers.AdminListProductOptions)
    admin.Post("/product_options", handlers.AdminCreateProductOption)
    admin.Put("/product_options/:id", handlers.AdminUpdateProductOption)
    admin.Delete("/product_options/:id", handlers.AdminDeleteProductOption)

    admin.Get("/orders", handlers.AdminListOrders)
    admin.Patch("/orders/:id/status", handlers.AdminUpdateOrderStatus)
}

