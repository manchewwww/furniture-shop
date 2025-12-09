package handlers

import (
    "github.com/gofiber/fiber/v2"

    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
)

// Departments
func AdminListDepartments(c *fiber.Ctx) error {
    var items []models.Department
    database.DB.Find(&items)
    return c.JSON(items)
}
func AdminCreateDepartment(c *fiber.Ctx) error {
    var in models.Department
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    if in.Name == "" { return c.Status(400).JSON(fiber.Map{"message":"Името е задължително"}) }
    if err := database.DB.Create(&in).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(in)
}
func AdminUpdateDepartment(c *fiber.Ctx) error {
    var in models.Department
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    in.ID = 0
    if err := database.DB.Model(&models.Department{}).Where("id = ?", c.Params("id")).Updates(in).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(fiber.Map{"message":"Обновено"})
}
func AdminDeleteDepartment(c *fiber.Ctx) error {
    database.DB.Delete(&models.Department{}, c.Params("id"))
    return c.JSON(fiber.Map{"message":"Изтрито"})
}

// Categories
func AdminListCategories(c *fiber.Ctx) error {
    var items []models.Category
    database.DB.Find(&items)
    return c.JSON(items)
}
func AdminCreateCategory(c *fiber.Ctx) error {
    var in models.Category
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    if in.Name == "" || in.DepartmentID == 0 { return c.Status(400).JSON(fiber.Map{"message":"Име и отдел са задължителни"}) }
    if err := database.DB.Create(&in).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(in)
}
func AdminUpdateCategory(c *fiber.Ctx) error {
    var in models.Category
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    in.ID = 0
    if err := database.DB.Model(&models.Category{}).Where("id = ?", c.Params("id")).Updates(in).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(fiber.Map{"message":"Обновено"})
}
func AdminDeleteCategory(c *fiber.Ctx) error {
    database.DB.Delete(&models.Category{}, c.Params("id"))
    return c.JSON(fiber.Map{"message":"Изтрито"})
}

// Products
func AdminListProducts(c *fiber.Ctx) error {
    var items []models.Product
    database.DB.Find(&items)
    return c.JSON(items)
}
func AdminCreateProduct(c *fiber.Ctx) error {
    var in models.Product
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    if in.Name == "" || in.CategoryID == 0 { return c.Status(400).JSON(fiber.Map{"message":"Име и категория са задължителни"}) }
    if err := database.DB.Create(&in).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(in)
}
func AdminUpdateProduct(c *fiber.Ctx) error {
    var in models.Product
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    in.ID = 0
    if err := database.DB.Model(&models.Product{}).Where("id = ?", c.Params("id")).Updates(in).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(fiber.Map{"message":"Обновено"})
}
func AdminDeleteProduct(c *fiber.Ctx) error {
    database.DB.Delete(&models.Product{}, c.Params("id"))
    return c.JSON(fiber.Map{"message":"Изтрито"})
}

// Product Options
func AdminListProductOptions(c *fiber.Ctx) error {
    var items []models.ProductOption
    pid := c.Query("product_id")
    if pid != "" {
        database.DB.Where("product_id = ?", pid).Find(&items)
    } else {
        database.DB.Find(&items)
    }
    return c.JSON(items)
}
func AdminCreateProductOption(c *fiber.Ctx) error {
    var in models.ProductOption
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    if in.ProductID == 0 || in.OptionType == "" || in.OptionName == "" { return c.Status(400).JSON(fiber.Map{"message":"Задължителни полета липсват"}) }
    if err := database.DB.Create(&in).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(in)
}
func AdminUpdateProductOption(c *fiber.Ctx) error {
    var in models.ProductOption
    if err := c.BodyParser(&in); err != nil { return c.Status(400).JSON(fiber.Map{"message":"Невалидни данни"}) }
    in.ID = 0
    if err := database.DB.Model(&models.ProductOption{}).Where("id = ?", c.Params("id")).Updates(in).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(fiber.Map{"message":"Обновено"})
}
func AdminDeleteProductOption(c *fiber.Ctx) error {
    database.DB.Delete(&models.ProductOption{}, c.Params("id"))
    return c.JSON(fiber.Map{"message":"Изтрито"})
}

