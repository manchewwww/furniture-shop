package handlers

import (
    "strconv"
    "github.com/gofiber/fiber/v2"

    models "furniture-shop/internal/domain/entity"
    app "furniture-shop/internal/app"
)

type CatalogHandler struct {
    svc app.CatalogService
}

func NewCatalogHandler(svc app.CatalogService) *CatalogHandler { return &CatalogHandler{svc: svc} }

func (h *CatalogHandler) GetDepartments() fiber.Handler {
    return func(c *fiber.Ctx) error {
        depts, err := h.svc.ListDepartments(c.Context())
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(depts)
    }
}

func (h *CatalogHandler) GetCategoriesByDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
        cats, err := h.svc.ListCategoriesByDepartment(c.Context(), uint(id))
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(cats)
    }
}

func (h *CatalogHandler) GetProductsByCategory() fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
        products, err := h.svc.ListProductsByCategory(c.Context(), uint(id))
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(products)
    }
}

func (h *CatalogHandler) GetProductDetails() fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
        p, err := h.svc.GetProduct(c.Context(), uint(id))
        if err != nil { return c.Status(404).JSON(fiber.Map{"message":"not found"}) }
        return c.JSON(p)
    }
}

func (h *CatalogHandler) SearchProducts() fiber.Handler {
    return func(c *fiber.Ctx) error {
        q := c.Query("query")
        if q == "" { return c.JSON([]models.Product{}) }
        items, err := h.svc.SearchProducts(c.Context(), q, 50)
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(items)
    }
}

func (h *CatalogHandler) GetProductRecommendations() fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
        rec, err := h.svc.RecommendProducts(c.Context(), uint(id), 4)
        if err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(rec)
    }
}




