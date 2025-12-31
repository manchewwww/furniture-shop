package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	ec "furniture-shop/internal/entities/catalog"
	ei "furniture-shop/internal/entities/inventory"
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
	var allCats []ec.Category
	if err := DB.Find(&allCats).Error; err != nil {
		return err
	}
	for _, cat := range allCats {
		prodCount := 8 + rand.Intn(3)
		for i := 1; i <= prodCount; i++ {
			p := ec.Product{
				CategoryID:             cat.ID,
				Name:                   cat.Name + " Product " + itoa(i),
				ShortDescription:       "Sample product for " + cat.Name,
				LongDescription:        "Auto-generated seed product.",
				BasePrice:              float64(100 + rand.Intn(900)),
				BaseProductionTimeDays: 7 + rand.Intn(21),
				ImageURL:               "https://via.placeholder.com/400x300",
				BaseMaterial:           []string{"MDF", "Wood", "Metal"}[rand.Intn(3)],
				DefaultWidth:           80 + rand.Intn(120),
				DefaultHeight:          30 + rand.Intn(170),
				DefaultDepth:           30 + rand.Intn(70),
			}
			if err := DB.Create(&p).Error; err != nil {
				return err
			}
		}
	}

	var some []ec.Product
	if err := DB.Limit(3).Find(&some).Error; err == nil && len(some) > 0 {
		opts := []ec.ProductOption{
			{ProductID: some[0].ID, OptionType: "material", OptionName: "Solid Wood", PriceModifierType: "percent", PriceModifierValue: 25, ProductionTimeModifierDays: 3},
			{ProductID: some[0].ID, OptionType: "extra", OptionName: "LED Lighting", PriceModifierType: "absolute", PriceModifierValue: 90, ProductionTimeModifierDays: 2},
		}
		for i := range opts {
			if err := DB.Create(&opts[i]).Error; err != nil {
				return err
			}
		}
	}

	admin := models.User{Role: "admin", Name: "Administrator", Email: "admin@example.com", Address: "Sofia", Phone: "+359888000000"}
	if err := admin.SetPassword("admin123"); err != nil {
		return err
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}
	upsertStock := func(material, unit string, qty float64) error {
		var s ei.Stock
		if err := DB.Where("material_name = ?", material).First(&s).Error; err != nil {
			if err := DB.Create(&ei.Stock{MaterialName: material, Unit: unit, QuantityAvailable: qty}).Error; err != nil {
				return err
			}
		}
		return nil
	}
	if err := upsertStock("MDF", "pcs", 1000); err != nil {
		return err
	}
	if err := upsertStock("Wood", "pcs", 1000); err != nil {
		return err
	}
	if err := upsertStock("Metal", "pcs", 1000); err != nil {
		return err
	}
	return nil
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }
