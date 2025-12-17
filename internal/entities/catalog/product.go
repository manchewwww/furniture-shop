package catalog

import "time"

type Product struct {
    ID                      uint            `gorm:"primaryKey" json:"id"`
    CategoryID              uint            `json:"category_id"`
    Name                    string          `json:"name"`
    ShortDescription        string          `json:"short_description"`
    LongDescription         string          `json:"long_description"`
    BasePrice               float64         `json:"base_price"`
    BaseProductionTimeDays  int             `json:"base_production_time_days"`
    ImageURL                string          `json:"image_url"`
    BaseMaterial            string          `json:"base_material"`
    DefaultWidth            int             `json:"default_width"`
    DefaultHeight           int             `json:"default_height"`
    DefaultDepth            int             `json:"default_depth"`
    IsMadeToOrder           bool            `json:"is_made_to_order"`
    CreatedAt               time.Time       `json:"created_at"`
    UpdatedAt               time.Time       `json:"updated_at"`
    Options                 []ProductOption `json:"options"`
}

type ProductOption struct {
    ID                            uint      `gorm:"primaryKey" json:"id"`
    ProductID                     uint      `json:"product_id"`
    OptionType                    string    `json:"option_type"`
    OptionName                    string    `json:"option_name"`
    PriceModifierType             string    `json:"price_modifier_type"`
    PriceModifierValue            float64   `json:"price_modifier_value"`
    ProductionTimeModifierDays    int       `json:"production_time_modifier_days"`
    ProductionTimeModifierPercent *int      `json:"production_time_modifier_percent"`
    CreatedAt                     time.Time `json:"created_at"`
    UpdatedAt                     time.Time `json:"updated_at"`
}

type RecommendationCounter struct {
    ID        uint `gorm:"primaryKey"`
    ProductID uint `gorm:"uniqueIndex"`
    Count     int
}

