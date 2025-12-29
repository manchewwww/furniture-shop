package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	ec "furniture-shop/internal/entities/catalog"
	models "furniture-shop/internal/entities/user"
)

func seedData() error {
	var count int64
	if err := DB.Model(&ec.Department{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Println("Seeding sample data...")
	depts := []ec.Department{
		{Name: "Living Room", Description: "Furniture for living spaces", ImageURL: "https://via.placeholder.com/600x400?text=Living+Room"},
		{Name: "Bedroom", Description: "Furniture for bedrooms", ImageURL: "https://via.placeholder.com/600x400?text=Bedroom"},
		{Name: "Kitchen", Description: "Furniture for kitchens", ImageURL: "https://via.placeholder.com/600x400?text=Kitchen"},
	}
	for i := range depts {
		if err := DB.Create(&depts[i]).Error; err != nil {
			return err
		}
	}

	var dLiving, dBedroom, dKitchen ec.Department
	DB.Where("name = ?", "Living Room").First(&dLiving)
	DB.Where("name = ?", "Bedroom").First(&dBedroom)
	DB.Where("name = ?", "Kitchen").First(&dKitchen)

	cats := []ec.Category{
		{Name: "Shelves", Description: "Wall and standing shelves", DepartmentID: dLiving.ID},
		{Name: "Storage Cabinets", Description: "Cabinets and sideboards", DepartmentID: dLiving.ID},
		{Name: "Beds", Description: "Single and double beds", DepartmentID: dBedroom.ID},
		{Name: "Wardrobes", Description: "Sliding and hinged wardrobes", DepartmentID: dBedroom.ID},
		{Name: "Kitchen Cabinets", Description: "Base and wall units", DepartmentID: dKitchen.ID},
		{Name: "Tables", Description: "Dining tables", DepartmentID: dKitchen.ID},
	}
	for i := range cats {
		if err := DB.Create(&cats[i]).Error; err != nil {
			return err
		}
	}

	rand.Seed(time.Now().UnixNano())
	var catWardrobes, catShelves, catBeds ec.Category
	DB.Where("name = ?", "Wardrobes").First(&catWardrobes)
	DB.Where("name = ?", "Shelves").First(&catShelves)
	DB.Where("name = ?", "Beds").First(&catBeds)

	var products []ec.Product
	for i := 1; i <= 10; i++ {
		products = append(products, ec.Product{
			CategoryID:             catWardrobes.ID,
			Name:                   "Modular Wardrobe " + itoa(i),
			ShortDescription:       "Modular wardrobe with customizable layout",
			LongDescription:        "Configurable wardrobe with options for shelves, drawers, and lighting.",
			BasePrice:              float64(700 + i*20),
			BaseProductionTimeDays: 14 + (i%3)*7,
			ImageURL:               "https://via.placeholder.com/400x300",
			BaseMaterial:           "MDF",
			DefaultWidth:           120 + i*5, DefaultHeight: 200, DefaultDepth: 60,
			IsMadeToOrder: true,
		})
	}
	for i := 1; i <= 8; i++ {
		products = append(products, ec.Product{
			CategoryID:             catShelves.ID,
			Name:                   "Wall Shelf " + itoa(i),
			ShortDescription:       "Minimal wall shelf",
			LongDescription:        "Compact shelving for books and decor with concealed mounts.",
			BasePrice:              float64(120 + i*10),
			BaseProductionTimeDays: 7 + (i%2)*3,
			ImageURL:               "https://via.placeholder.com/400x300",
			BaseMaterial:           "MDF",
			DefaultWidth:           80 + i*5, DefaultHeight: 30, DefaultDepth: 25,
			IsMadeToOrder: false,
		})
	}
	for i := 1; i <= 9; i++ {
		products = append(products, ec.Product{
			CategoryID:             catBeds.ID,
			Name:                   "Double Bed " + itoa(i),
			ShortDescription:       "Comfortable bed",
			LongDescription:        "Sturdy bed frame with optional storage and headboard.",
			BasePrice:              float64(650 + i*30),
			BaseProductionTimeDays: 21,
			ImageURL:               "https://via.placeholder.com/400x300",
			BaseMaterial:           "Wood",
			DefaultWidth:           160, DefaultHeight: 40, DefaultDepth: 200,
			IsMadeToOrder: true,
		})
	}
	for i := range products {
		if err := DB.Create(&products[i]).Error; err != nil {
			return err
		}
	}

	var a, b, c ec.Product
	DB.Where("category_id = ?", catWardrobes.ID).First(&a)
	DB.Where("category_id = ?", catShelves.ID).First(&b)
	DB.Where("category_id = ?", catBeds.ID).First(&c)
	opts := []ec.ProductOption{
		{ProductID: a.ID, OptionType: "material", OptionName: "Solid Wood", PriceModifierType: "percent", PriceModifierValue: 25, ProductionTimeModifierDays: 3},
		{ProductID: a.ID, OptionType: "extra", OptionName: "LED Lighting", PriceModifierType: "absolute", PriceModifierValue: 90, ProductionTimeModifierDays: 2},
		{ProductID: c.ID, OptionType: "extra", OptionName: "Mattress", PriceModifierType: "absolute", PriceModifierValue: 120, ProductionTimeModifierDays: 4},
		// Color options
		{ProductID: a.ID, OptionType: "color", OptionName: "White", PriceModifierType: "absolute", PriceModifierValue: 0},
		{ProductID: a.ID, OptionType: "color", OptionName: "Oak", PriceModifierType: "absolute", PriceModifierValue: 0},
		{ProductID: a.ID, OptionType: "color", OptionName: "Walnut", PriceModifierType: "absolute", PriceModifierValue: 0},
	}
	for i := range opts {
		if err := DB.Create(&opts[i]).Error; err != nil {
			return err
		}
	}

	admin := models.User{Role: "admin", Name: "Administrator", Email: "admin@example.com", Address: "Sofia", Phone: "+359888000000"}
	if err := admin.SetPassword("admin123"); err != nil {
		return err
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}
	return nil
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }
