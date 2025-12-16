package models

type RecommendationCounter struct {
    ID        uint `gorm:"primaryKey"`
    ProductID uint `gorm:"uniqueIndex"`
    Count     int
}

