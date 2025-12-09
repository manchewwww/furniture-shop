package handlers

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"

    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
)

func GetDepartments(c *fiber.Ctx) error {
    var depts []models.Department
    if err := database.DB.Find(&depts).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(depts)
}

func GetCategoriesByDepartment(c *fiber.Ctx) error {
    var cats []models.Category
    if err := database.DB.Where("department_id = ?", c.Params("id")).Find(&cats).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(cats)
}

func GetProductsByCategory(c *fiber.Ctx) error {
    var products []models.Product
    if err := database.DB.Where("category_id = ?", c.Params("id")).Find(&products).Error; err != nil { return c.Status(500).JSON(fiber.Map{"message":"Грешка"}) }
    return c.JSON(products)
}

func GetProductDetails(c *fiber.Ctx) error {
    var p models.Product
    if err := database.DB.Preload("Options").First(&p, c.Params("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound { return c.Status(404).JSON(fiber.Map{"message":"Продуктът не е намерен"}) }
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(p)
}

func SearchProducts(c *fiber.Ctx) error {
    q := c.Query("query")
    var items []models.Product
    if q == "" { return c.JSON([]models.Product{}) }
    like := "%" + q + "%"
    if err := database.DB.Where("name ILIKE ? OR short_description ILIKE ?", like, like).Limit(50).Find(&items).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(items)
}

func GetProductRecommendations(c *fiber.Ctx) error {
    // simple: same category top 4 by id
    var p models.Product
    if err := database.DB.First(&p, c.Params("id")).Error; err != nil { return c.Status(404).JSON(fiber.Map{"message":"Не е намерен"}) }
    var rec []models.Product
    if err := database.DB.Where("category_id = ? AND id <> ?", p.CategoryID, p.ID).Limit(4).Find(&rec).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message":"Грешка"})
    }
    return c.JSON(rec)
}

