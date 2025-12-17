package database

import (
    "fmt"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "furniture-shop/internal/config"
    ec "furniture-shop/internal/entities/catalog"
    ei "furniture-shop/internal/entities/inventory"
    eo "furniture-shop/internal/entities/orders"
    eu "furniture-shop/internal/entities/user"
)

const DATABASE_PREFIX = "DATABASE"

var DB *gorm.DB

func Connect() error {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	if user == "" || password == "" {
		return fmt.Errorf("%s: missing database credentials", DATABASE_PREFIX)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.Configurations.DB.HOST, user, password, config.Configurations.DB.NAME,
		config.Configurations.DB.PORT, config.Configurations.DB.SSL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	return nil
}

func AutoMigrateAndSeed() error {
    if err := DB.AutoMigrate(
        &ec.Department{},
        &ec.Category{},
        &ec.Product{},
        &ec.ProductOption{},
        &eu.User{},
        &eo.Order{},
        &eo.OrderItem{},
        &ei.Stock{},
        &ec.RecommendationCounter{},
    ); err != nil {
        return err
    }
    return seedData()
}
