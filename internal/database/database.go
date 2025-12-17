package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"furniture-shop/internal/config"
	models "furniture-shop/internal/domain/entity"
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
		&models.Department{},
		&models.Category{},
		&models.Product{},
		&models.ProductOption{},
		&models.User{},
		&models.Order{},
		&models.OrderItem{},
		&models.Stock{},
		&models.RecommendationCounter{},
	); err != nil {
		return err
	}
	return seedData()
}
