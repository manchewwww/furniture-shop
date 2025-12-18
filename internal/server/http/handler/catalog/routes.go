package catalog

import "github.com/gofiber/fiber/v2"

func Register(api fiber.Router, h *Handler) {
    api.Get("/departments", h.GetDepartments())
    api.Get("/departments/:id/categories", h.GetCategoriesByDepartment())
    api.Get("/categories/:id/products", h.GetProductsByCategory())
    api.Get("/products/:id", h.GetProductDetails())
    api.Get("/products/:id/recommendations", h.GetProductRecommendations())
    api.Get("/products/search", h.SearchProducts())
}

