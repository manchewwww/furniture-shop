package orders

type CreateOrderItem struct {
	ProductID uint             `json:"product_id" validate:"required,gt=0"`
	Quantity  int              `json:"quantity" validate:"required,gt=0"`
	Options   []SelectedOption `json:"options"`
}
