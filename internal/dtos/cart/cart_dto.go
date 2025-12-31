package cart

type ReplaceCartRequest struct {
	Items []CartItemInput `json:"items"`
}

type AddCartItemRequest struct {
	ProductID uint             `json:"product_id"`
	Quantity  int              `json:"quantity"`
	Options   []SelectedOption `json:"options"`
}

type UpdateCartItemRequest struct {
	Quantity int              `json:"quantity"`
	Options  []SelectedOption `json:"options"`
}

type CartItemInput struct {
	ProductID uint             `json:"product_id"`
	Quantity  int              `json:"quantity"`
	Options   []SelectedOption `json:"options"`
}

type SelectedOption struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}
