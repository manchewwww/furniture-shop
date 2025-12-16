package handlers

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    models "furniture-shop/internal/domain/entity"
    "furniture-shop/internal/services"
)

type AdminHandler struct { svc services.AdminService }

func NewAdminHandler(svc services.AdminService) *AdminHandler { return &AdminHandler{svc: svc} }

// Departments
func (h *AdminHandler) ListDepartments() fiber.Handler {
    return func(c *fiber.Ctx) error { items, _ := h.svc.ListDepartments(c.Context()); return c.JSON(items) }
}
func (h *AdminHandler) CreateDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in models.Department
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        if in.Name == "" { return c.Status(400).JSON(fiber.Map{"message":"name is required"}) }
        if err := h.svc.CreateDepartment(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(in)
    }
}
func (h *AdminHandler) UpdateDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in models.Department
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        var id uint; _, _ = fmt.Sscan(c.Params("id"), &id)
        if err := h.svc.UpdateDepartment(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(fiber.Map{"message":"updated"})
    }
}
func (h *AdminHandler) DeleteDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteDepartment(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

// Categories
func (h *AdminHandler) ListCategories() fiber.Handler { return func(c *fiber.Ctx) error { items, _ := h.svc.ListCategories(c.Context()); return c.JSON(items) } }
func (h *AdminHandler) CreateCategory() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in models.Category
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        if in.Name == "" || in.DepartmentID == 0 { return c.Status(400).JSON(fiber.Map{"message":"name and department_id are required"}) }
        if err := h.svc.CreateCategory(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(in)
    }
}
func (h *AdminHandler) UpdateCategory() fiber.Handler {
    return func(c *fiber.Ctx) error { var in models.Category; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); if err := h.svc.UpdateCategory(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(fiber.Map{"message":"updated"}) }
}
func (h *AdminHandler) DeleteCategory() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteCategory(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

// Products
func (h *AdminHandler) ListProducts() fiber.Handler { return func(c *fiber.Ctx) error { items, _ := h.svc.ListProducts(c.Context()); return c.JSON(items) } }
func (h *AdminHandler) CreateProduct() fiber.Handler {
    return func(c *fiber.Ctx) error { var in models.Product; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; if in.Name == "" || in.CategoryID == 0 { return c.Status(400).JSON(fiber.Map{"message":"name and category_id are required"}) }; if err := h.svc.CreateProduct(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(in) }
}
func (h *AdminHandler) UpdateProduct() fiber.Handler {
    return func(c *fiber.Ctx) error { var in models.Product; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); if err := h.svc.UpdateProduct(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(fiber.Map{"message":"updated"}) }
}
func (h *AdminHandler) DeleteProduct() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteProduct(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

// Product Options
func (h *AdminHandler) ListProductOptions() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var pid *uint
        if s := c.Query("product_id"); s != "" { var v uint; _, _ = fmt.Sscan(s, &v); pid = &v }
        items, _ := h.svc.ListProductOptions(c.Context(), pid)
        return c.JSON(items)
    }
}
func (h *AdminHandler) CreateProductOption() fiber.Handler {
    return func(c *fiber.Ctx) error { var in models.ProductOption; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; if in.ProductID == 0 || in.OptionType == "" || in.OptionName == "" { return c.Status(400).JSON(fiber.Map{"message":"missing fields"}) }; if err := h.svc.CreateProductOption(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(in) }
}
func (h *AdminHandler) UpdateProductOption() fiber.Handler {
    return func(c *fiber.Ctx) error { var in models.ProductOption; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); if err := h.svc.UpdateProductOption(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(fiber.Map{"message":"updated"}) }
}
func (h *AdminHandler) DeleteProductOption() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteProductOption(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}
