package admin

import (
    "fmt"
    "github.com/gofiber/fiber/v2"

    ec "furniture-shop/internal/entities/catalog"
    "furniture-shop/internal/service"
)

type Handler struct { svc service.AdminService }

func NewAdminHandler(svc service.AdminService) *Handler { return &Handler{svc: svc} }

// Departments
func (h *Handler) ListDepartments() fiber.Handler { return func(c *fiber.Ctx) error { items, _ := h.svc.ListDepartments(c.Context()); return c.JSON(items) } }
func (h *Handler) CreateDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in ec.Department
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        if in.Name == "" { return c.Status(400).JSON(fiber.Map{"message":"name is required"}) }
        if err := h.svc.CreateDepartment(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(in)
    }
}
func (h *Handler) UpdateDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in ec.Department
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        var id uint; _, _ = fmt.Sscan(c.Params("id"), &id)
        if err := h.svc.UpdateDepartment(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(fiber.Map{"message":"updated"})
    }
}
func (h *Handler) DeleteDepartment() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteDepartment(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

// Categories
func (h *Handler) ListCategories() fiber.Handler { return func(c *fiber.Ctx) error { items, _ := h.svc.ListCategories(c.Context()); return c.JSON(items) } }
func (h *Handler) CreateCategory() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var in ec.Category
        if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }
        if in.Name == "" || in.DepartmentID == 0 { return c.Status(400).JSON(fiber.Map{"message":"name and department_id are required"}) }
        if err := h.svc.CreateCategory(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }
        return c.JSON(in)
    }
}
func (h *Handler) UpdateCategory() fiber.Handler {
    return func(c *fiber.Ctx) error { var in ec.Category; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); if err := h.svc.UpdateCategory(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(fiber.Map{"message":"updated"}) }
}
func (h *Handler) DeleteCategory() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteCategory(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

// Products
func (h *Handler) ListProducts() fiber.Handler { return func(c *fiber.Ctx) error { items, _ := h.svc.ListProducts(c.Context()); return c.JSON(items) } }
func (h *Handler) CreateProduct() fiber.Handler {
    return func(c *fiber.Ctx) error { var in ec.Product; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; if in.Name == "" || in.CategoryID == 0 { return c.Status(400).JSON(fiber.Map{"message":"name and category_id are required"}) }; if err := h.svc.CreateProduct(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(in) }
}
func (h *Handler) UpdateProduct() fiber.Handler {
    return func(c *fiber.Ctx) error { var in ec.Product; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); if err := h.svc.UpdateProduct(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(fiber.Map{"message":"updated"}) }
}
func (h *Handler) DeleteProduct() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteProduct(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

// Product Options
func (h *Handler) ListProductOptions() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var pid *uint
        if s := c.Query("product_id"); s != "" { var v uint; _, _ = fmt.Sscan(s, &v); pid = &v }
        items, _ := h.svc.ListProductOptions(c.Context(), pid)
        return c.JSON(items)
    }
}
func (h *Handler) CreateProductOption() fiber.Handler {
    return func(c *fiber.Ctx) error { var in ec.ProductOption; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; if in.ProductID == 0 || in.OptionType == "" || in.OptionName == "" { return c.Status(400).JSON(fiber.Map{"message":"missing fields"}) }; if err := h.svc.CreateProductOption(c.Context(), &in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(in) }
}
func (h *Handler) UpdateProductOption() fiber.Handler {
    return func(c *fiber.Ctx) error { var in ec.ProductOption; if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"invalid request"}) }; var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); if err := h.svc.UpdateProductOption(c.Context(), id, in); err != nil { return c.Status(500).JSON(fiber.Map{"message":"server error"}) }; return c.JSON(fiber.Map{"message":"updated"}) }
}
func (h *Handler) DeleteProductOption() fiber.Handler {
    return func(c *fiber.Ctx) error { var id uint; _, _ = fmt.Sscan(c.Params("id"), &id); _ = h.svc.DeleteProductOption(c.Context(), id); return c.JSON(fiber.Map{"message":"deleted"}) }
}

