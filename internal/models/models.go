package models

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type Department struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    DepartmentID uint      `json:"department_id"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type Product struct {
    ID                      uint      `gorm:"primaryKey" json:"id"`
    CategoryID              uint      `json:"category_id"`
    Name                    string    `json:"name"`
    ShortDescription        string    `json:"short_description"`
    LongDescription         string    `json:"long_description"`
    BasePrice               float64   `json:"base_price"`
    BaseProductionTimeDays  int       `json:"base_production_time_days"`
    ImageURL                string    `json:"image_url"`
    BaseMaterial            string    `json:"base_material"`
    DefaultWidth            int       `json:"default_width"`
    DefaultHeight           int       `json:"default_height"`
    DefaultDepth            int       `json:"default_depth"`
    IsMadeToOrder           bool      `json:"is_made_to_order"`
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
    Options                 []ProductOption `json:"options"`
}

type ProductOption struct {
    ID                           uint      `gorm:"primaryKey" json:"id"`
    ProductID                    uint      `json:"product_id"`
    OptionType                   string    `json:"option_type"`
    OptionName                   string    `json:"option_name"`
    PriceModifierType            string    `json:"price_modifier_type"`
    PriceModifierValue           float64   `json:"price_modifier_value"`
    ProductionTimeModifierDays   int       `json:"production_time_modifier_days"`
    ProductionTimeModifierPercent *int     `json:"production_time_modifier_percent"`
    CreatedAt                    time.Time `json:"created_at"`
    UpdatedAt                    time.Time `json:"updated_at"`
}

type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Role         string    `json:"role"`
    Name         string    `json:"name"`
    Email        string    `gorm:"uniqueIndex" json:"email"`
    PasswordHash string    `json:"-"`
    Address      string    `json:"address"`
    Phone        string    `json:"phone"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) SetPassword(plain string) error {
    h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil { return err }
    u.PasswordHash = string(h)
    return nil
}

func (u *User) CheckPassword(plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain)) == nil
}

type Order struct {
    ID                          uint        `gorm:"primaryKey" json:"id"`
    UserID                      uint        `json:"user_id"`
    Status                      string      `json:"status"`
    TotalPrice                  float64     `json:"total_price"`
    EstimatedProductionTimeDays int         `json:"estimated_production_time_days"`
    PaymentMethod               string      `json:"payment_method"`
    PaymentStatus               string      `json:"payment_status"`
    CreatedAt                   time.Time   `json:"created_at"`
    UpdatedAt                   time.Time   `json:"updated_at"`
    Items                       []OrderItem `json:"items"`
}

type OrderItem struct {
    ID                           uint      `gorm:"primaryKey" json:"id"`
    OrderID                      uint      `json:"order_id"`
    ProductID                    uint      `json:"product_id"`
    Quantity                     int       `json:"quantity"`
    UnitPrice                    float64   `json:"unit_price"`
    LineTotal                    float64   `json:"line_total"`
    CalculatedProductionTimeDays int       `json:"calculated_production_time_days"`
    SelectedOptionsJSON          string    `json:"selected_options_json"`
    CreatedAt                    time.Time `json:"created_at"`
    UpdatedAt                    time.Time `json:"updated_at"`
}

type Stock struct {
    ID                uint      `gorm:"primaryKey" json:"id"`
    MaterialName      string    `json:"material_name"`
    QuantityAvailable float64   `json:"quantity_available"`
    Unit              string    `json:"unit"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

type RecommendationCounter struct {
    ID        uint `gorm:"primaryKey"`
    ProductID uint `gorm:"uniqueIndex"`
    Count     int
}

