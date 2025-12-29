package orders

import "time"

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
