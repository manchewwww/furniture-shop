package orders

type CreateOrderInput struct {
	UserID        *uint             `json:"user_id"`
	Name          string            `json:"name" validate:"required,min=2"`
	Email         string            `json:"email" validate:"required,email"`
	Address       string            `json:"address" validate:"required,min=5"`
	Phone         string            `json:"phone" validate:"required,phone"`
	Items         []CreateOrderItem `json:"items" validate:"required,min=1,dive"`
	PaymentMethod string            `json:"payment_method" validate:"required,oneof=card"`
}
