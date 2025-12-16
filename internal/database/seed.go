package database

import (
    "fmt"
    "log"
    "math/rand"
    "time"

    "furniture-shop/internal/models"
)

func seedData() error {
    var count int64
    if err := DB.Model(&models.Department{}).Count(&count).Error; err != nil { return err }
    if count > 0 { return nil }

    log.Println("Seeding sample data...")
    // Departments
    depts := []models.Department{
        {Name: "Дневна", Description: "Мебели за дневна стая"},
        {Name: "Спалня", Description: "Мебели за спалня"},
        {Name: "Кухня", Description: "Мебели за кухня"},
    }
    for i := range depts { if err := DB.Create(&depts[i]).Error; err != nil { return err } }

    var dLiving, dBedroom, dKitchen models.Department
    DB.Where("name = ?", "Дневна").First(&dLiving)
    DB.Where("name = ?", "Спалня").First(&dBedroom)
    DB.Where("name = ?", "Кухня").First(&dKitchen)

    // Categories per department
    cats := []models.Category{
        {Name: "Шкафове", Description: "Шкафове и витрини", DepartmentID: dLiving.ID},
        {Name: "Етажерки", Description: "Етажерки и рафтове", DepartmentID: dLiving.ID},
        {Name: "Легла", Description: "Легла и нощни шкафчета", DepartmentID: dBedroom.ID},
        {Name: "Гардероби", Description: "Гардероби", DepartmentID: dBedroom.ID},
        {Name: "Кухненски шкафове", Description: "Горни и долни модули", DepartmentID: dKitchen.ID},
        {Name: "Маси", Description: "Кухненски маси", DepartmentID: dKitchen.ID},
    }
    for i := range cats { if err := DB.Create(&cats[i]).Error; err != nil { return err } }

    // Products: add 8-10 per some categories
    rand.Seed(time.Now().UnixNano())
    var catWardrobes, catShelves, catBeds models.Category
    DB.Where("name = ?", "Гардероби").First(&catWardrobes)
    DB.Where("name = ?", "Етажерки").First(&catShelves)
    DB.Where("name = ?", "Легла").First(&catBeds)

    var products []models.Product
    for i := 1; i <= 10; i++ {
        products = append(products, models.Product{
            CategoryID: catWardrobes.ID,
            Name:        "Гардероб Модел " + itoa(i),
            ShortDescription: "Гардероб с плъзгащи врати",
            LongDescription:  "Качествен гардероб с възможност за избор на материал и размери.",
            BasePrice:  float64(700 + i*20),
            BaseProductionTimeDays: 14 + (i%3)*7,
            ImageURL:   "https://via.placeholder.com/400x300",
            BaseMaterial: "MDF",
            DefaultWidth: 120 + i*5, DefaultHeight: 200, DefaultDepth: 60,
            IsMadeToOrder: true,
        })
    }
    for i := 1; i <= 8; i++ {
        products = append(products, models.Product{
            CategoryID: catShelves.ID,
            Name:        "Етажерка Модел " + itoa(i),
            ShortDescription: "Стенна етажерка",
            LongDescription:  "Практична етажерка за книги и декорации.",
            BasePrice:  float64(120 + i*10),
            BaseProductionTimeDays: 7 + (i%2)*3,
            ImageURL:   "https://via.placeholder.com/400x300",
            BaseMaterial: "MDF",
            DefaultWidth: 80 + i*5, DefaultHeight: 30, DefaultDepth: 25,
            IsMadeToOrder: false,
        })
    }
    for i := 1; i <= 9; i++ {
        products = append(products, models.Product{
            CategoryID: catBeds.ID,
            Name:        "Легло Комфорт " + itoa(i),
            ShortDescription: "Двойно легло",
            LongDescription:  "Удобно двойно легло с опция за чекмеджета.",
            BasePrice:  float64(650 + i*30),
            BaseProductionTimeDays: 21,
            ImageURL:   "https://via.placeholder.com/400x300",
            BaseMaterial: "масив",
            DefaultWidth: 160, DefaultHeight: 40, DefaultDepth: 200,
            IsMadeToOrder: true,
        })
    }
    for i := range products { if err := DB.Create(&products[i]).Error; err != nil { return err } }

    // Add a few options to first product of each category
    var a, b, c models.Product
    DB.Where("category_id = ?", catWardrobes.ID).First(&a)
    DB.Where("category_id = ?", catShelves.ID).First(&b)
    DB.Where("category_id = ?", catBeds.ID).First(&c)
    opts := []models.ProductOption{
        {ProductID: a.ID, OptionType: "material", OptionName: "Масив", PriceModifierType: "percent", PriceModifierValue: 25, ProductionTimeModifierDays: 3},
        {ProductID: a.ID, OptionType: "extra", OptionName: "LED осветление", PriceModifierType: "absolute", PriceModifierValue: 90, ProductionTimeModifierDays: 2},
        {ProductID: c.ID, OptionType: "extra", OptionName: "Чекмеджета", PriceModifierType: "absolute", PriceModifierValue: 120, ProductionTimeModifierDays: 4},
    }
    for i := range opts { if err := DB.Create(&opts[i]).Error; err != nil { return err } }

    // Admin user (password: admin123)
    admin := models.User{Role: "admin", Name: "Администратор", Email: "admin@example.com", Address: "София", Phone: "+359888000000"}
    if err := admin.SetPassword("admin123"); err != nil { return err }
    if err := DB.Create(&admin).Error; err != nil { return err }
    return nil
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }
