package database

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	ec "furniture-shop/internal/entities/catalog"
	models "furniture-shop/internal/entities/user"
)

func seedData() error {
	if strings.EqualFold(os.Getenv("SEED_RESET"), "true") {
		_ = DB.Exec("TRUNCATE TABLE product_options, products, categories, departments, recommendation_counters, users RESTART IDENTITY CASCADE").Error
	}
	var count int64
	if err := DB.Model(&ec.Department{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Println("Seeding furniture shop data (English)...")
	depts := []ec.Department{
		{Name: "Living Room", Description: "Sofas, armchairs, coffee tables, TV units", ImageURL: ""},
		{Name: "Bedroom", Description: "Beds, wardrobes, dressers, nightstands", ImageURL: ""},
		{Name: "Dining", Description: "Dining tables, chairs, sideboards", ImageURL: ""},
		{Name: "Home Office", Description: "Desks, office chairs, bookcases", ImageURL: ""},
		{Name: "Outdoor", Description: "Garden tables, chairs, loungers", ImageURL: ""},
	}
	for i := range depts {
		if err := DB.Create(&depts[i]).Error; err != nil {
			return err
		}
	}

	var dLiving, dBedroom, dDining, dOffice, dOutdoor ec.Department
	DB.Where("name = ?", "Living Room").First(&dLiving)
	DB.Where("name = ?", "Bedroom").First(&dBedroom)
	DB.Where("name = ?", "Dining").First(&dDining)
	DB.Where("name = ?", "Home Office").First(&dOffice)
	DB.Where("name = ?", "Outdoor").First(&dOutdoor)

	cats := []ec.Category{
		// Living Room (5)
		{Name: "Sofas", Description: "Two- and three-seater sofas", DepartmentID: dLiving.ID},
		{Name: "Armchairs", Description: "Comfortable lounge chairs", DepartmentID: dLiving.ID},
		{Name: "Coffee Tables", Description: "Coffee & side tables", DepartmentID: dLiving.ID},
		{Name: "TV Units", Description: "TV stands and media units", DepartmentID: dLiving.ID},
		{Name: "Console Tables", Description: "Hall and console tables", DepartmentID: dLiving.ID},

		// Bedroom (5)
		{Name: "Beds", Description: "Beds and headboards", DepartmentID: dBedroom.ID},
		{Name: "Wardrobes", Description: "Hinged and sliding wardrobes", DepartmentID: dBedroom.ID},
		{Name: "Dressers", Description: "Dressers and chests of drawers", DepartmentID: dBedroom.ID},
		{Name: "Nightstands", Description: "Bedside tables", DepartmentID: dBedroom.ID},
		{Name: "Mattresses", Description: "Mattresses and toppers", DepartmentID: dBedroom.ID},

		// Dining (5)
		{Name: "Dining Tables", Description: "Tables for every space", DepartmentID: dDining.ID},
		{Name: "Dining Chairs", Description: "Upholstered and wooden chairs", DepartmentID: dDining.ID},
		{Name: "Sideboards", Description: "Storage for dining rooms", DepartmentID: dDining.ID},
		{Name: "Bar Stools", Description: "Stools for bars and islands", DepartmentID: dDining.ID},
		{Name: "Dining Benches", Description: "Benches for tables", DepartmentID: dDining.ID},

		// Home Office (5)
		{Name: "Desks", Description: "Work and study desks", DepartmentID: dOffice.ID},
		{Name: "Office Chairs", Description: "Ergonomic seating", DepartmentID: dOffice.ID},
		{Name: "Bookcases", Description: "Bookcases and shelving", DepartmentID: dOffice.ID},
		{Name: "Filing Cabinets", Description: "Document storage", DepartmentID: dOffice.ID},
		{Name: "Shelving Units", Description: "Open shelving", DepartmentID: dOffice.ID},

		// Outdoor (5)
		{Name: "Outdoor Tables", Description: "Garden and patio tables", DepartmentID: dOutdoor.ID},
		{Name: "Outdoor Chairs", Description: "Chairs and loungers", DepartmentID: dOutdoor.ID},
		{Name: "Loungers", Description: "Sun loungers", DepartmentID: dOutdoor.ID},
		{Name: "Outdoor Sofas", Description: "Garden sofas and sets", DepartmentID: dOutdoor.ID},
		{Name: "Parasols", Description: "Umbrellas and shades", DepartmentID: dOutdoor.ID},
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
	cityNames := []string{"Sofia", "Plovdiv", "Varna", "Burgas", "Ruse", "Stara Zagora", "Veliko Tarnovo", "Blagoevgrad"}
	mats := []string{"MDF", "Wood", "Metal"}
	for _, cat := range allCats {
		for i := 1; i <= 10; i++ { // exactly 10 products per category
			city := cityNames[rand.Intn(len(cityNames))]
			name := fmt.Sprintf("%s %s %d", city, cat.Name, i)
			p := ec.Product{
				CategoryID:             cat.ID,
				Name:                   name,
				ShortDescription:       fmt.Sprintf("%s for modern homes", cat.Name),
				LongDescription:        fmt.Sprintf("Quality %s crafted with durable finishes and clean lines.", cat.Name),
				BasePrice:              float64(150 + rand.Intn(1200)),
				BaseProductionTimeDays: 7 + rand.Intn(21),
				ImageURL:               findUploadImage(cat.Name, i),
				BaseMaterial:           mats[rand.Intn(len(mats))],
				DefaultWidth:           60 + rand.Intn(140),
				DefaultHeight:          35 + rand.Intn(180),
				DefaultDepth:           30 + rand.Intn(70),
			}
			if err := DB.Create(&p).Error; err != nil {
				return err
			}
		}
	}

	// helper to make URL-friendly file names for categories
	// e.g. "Coffee Tables" -> "coffee-tables"

	// Options for first 10 products
	var some []ec.Product
	if err := DB.Limit(10).Find(&some).Error; err == nil && len(some) > 0 {
		for _, pr := range some {
			opts := []ec.ProductOption{
				{ProductID: pr.ID, OptionType: "color", OptionName: "White", PriceModifierType: "absolute", PriceModifierValue: 0},
				{ProductID: pr.ID, OptionType: "color", OptionName: "Oak", PriceModifierType: "absolute", PriceModifierValue: 0},
				{ProductID: pr.ID, OptionType: "color", OptionName: "Walnut", PriceModifierType: "absolute", PriceModifierValue: 0},
				{ProductID: pr.ID, OptionType: "material", OptionName: "Solid Wood", PriceModifierType: "percent", PriceModifierValue: 20, ProductionTimeModifierDays: 3},
				{ProductID: pr.ID, OptionType: "extra", OptionName: "Soft-Close Hinges", PriceModifierType: "absolute", PriceModifierValue: 35, ProductionTimeModifierDays: 1},
			}
			for i := range opts {
				_ = DB.Create(&opts[i]).Error
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
	return nil
}

func slug(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			out = append(out, r+('a'-'A'))
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			out = append(out, r)
		case r == ' ' || r == '_' || r == '-' || r == '/' || r == '+':
			out = append(out, '-')
		default:
			// skip other symbols
		}
	}
	res := make([]rune, 0, len(out))
	var prevDash bool
	for _, r := range out {
		if r == '-' {
			if !prevDash {
				res = append(res, r)
				prevDash = true
			}
		} else {
			res = append(res, r)
			prevDash = false
		}
	}
	return string(res)
}

func findUploadImage(categoryName string, idx int) string {
	baseSlug := slug(categoryName)
	baseName := strings.TrimSpace(categoryName)
	try := func(fn string) string {
		p := filepath.Join("uploads", fn)
		if _, err := os.Stat(p); err == nil {
			return "/uploads/" + fn
		}
		return ""
	}
	exts := []string{"jpg", "jpeg", "png", "webp", "avif"}

	aliasCandidates := []string{baseSlug}
	rawLower := strings.ToLower(baseName)
	aliasCandidates = append(aliasCandidates, strings.ReplaceAll(rawLower, " ", "-"))
	aliasCandidates = append(aliasCandidates, rawLower)
	if strings.HasSuffix(baseSlug, "s") {
		aliasCandidates = append(aliasCandidates, strings.TrimSuffix(baseSlug, "s"))
	}
	if strings.HasSuffix(rawLower, "s") {
		aliasCandidates = append(aliasCandidates, strings.TrimSuffix(strings.ReplaceAll(rawLower, " ", "-"), "s"))
		aliasCandidates = append(aliasCandidates, strings.TrimSuffix(rawLower, "s"))
	}
	irregular := map[string][]string{
		"mattresses": {"matters"},
		"wardrobes":  {"wardrob", "wardrobe"},
		"sofas":      {"sofa"},
		"beds":       {"bed"},
	}
	if v, ok := irregular[strings.ToLower(baseName)]; ok {
		aliasCandidates = append(aliasCandidates, v...)
	}

	for _, base := range aliasCandidates {
		for _, ext := range exts {
			if r := try(fmt.Sprintf("%s-%d.%s", base, idx, ext)); r != "" {
				return r
			}
			if r := try(fmt.Sprintf("%s_%d.%s", base, idx, ext)); r != "" {
				return r
			}
			if r := try(fmt.Sprintf("%s %d.%s", base, idx, ext)); r != "" {
				return r
			}
		}
	}

	files, _ := os.ReadDir("uploads")
	if len(files) > 0 {
		norm := func(s string) string {
			s = strings.ToLower(s)
			b := strings.Builder{}
			for _, r := range s {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == ' ' || r == '_' {
					b.WriteRune(r)
				}
			}
			return b.String()
		}
		want := norm(baseSlug)
		wantSing := strings.TrimSuffix(want, "s")
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			lower := strings.ToLower(name)
			for _, ext := range exts {
				suffixes := []string{fmt.Sprintf("-%d.%s", idx, ext), fmt.Sprintf("_%d.%s", idx, ext), fmt.Sprintf(" %d.%s", idx, ext)}
				for _, suf := range suffixes {
					if strings.HasSuffix(lower, suf) {
						base := strings.TrimSuffix(lower, suf)
						n := norm(base)
						if strings.Contains(n, want) || strings.Contains(n, wantSing) || strings.HasPrefix(want, n) || strings.HasPrefix(n, want) {
							return "/uploads/" + name
						}
					}
				}
			}
		}
	}
	return ""
}
