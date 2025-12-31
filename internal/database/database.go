package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"furniture-shop/internal/config"
	ec "furniture-shop/internal/entities/catalog"
	eo "furniture-shop/internal/entities/orders"
	eu "furniture-shop/internal/entities/user"
)

const databaseErrorPrefix = "DATABASE"

var DB *gorm.DB

func Connect() error {
	user := config.Env.DBUser
	password := config.Env.DBPass
	if user == "" || password == "" {
		return fmt.Errorf("%s: missing database credentials (DB_USER and DB_PASSWORD environment variables)", databaseErrorPrefix)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.Configurations.DB.Host, user, password, config.Configurations.DB.Name,
		config.Configurations.DB.Port, config.Configurations.DB.SSL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(20)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(60 * time.Minute)
	}

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
		&eo.Cart{},
		&eo.CartItem{},
		&ec.RecommendationCounter{},
	); err != nil {
		return err
	}
	return seedData()
}
