package orders

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"uniqueIndex" json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	CartID              uint      `json:"cart_id"`
	ProductID           uint      `json:"product_id"`
	Quantity            int       `json:"quantity"`
	SelectedOptionsJSON string    `json:"selected_options_json"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
