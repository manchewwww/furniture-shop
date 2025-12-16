package database

import (
    "fmt"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "furniture-shop/internal/config"
    "furniture-shop/internal/models"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Sofia",
        cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetMaxOpenConns(10)
    sqlDB.SetConnMaxLifetime(time.Hour)
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

