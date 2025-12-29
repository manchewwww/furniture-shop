package inventory

import "time"

type Stock struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	MaterialName      string    `json:"material_name"`
	QuantityAvailable float64   `json:"quantity_available"`
	Unit              string    `json:"unit"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
