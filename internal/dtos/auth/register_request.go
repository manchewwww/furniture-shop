package auth

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Address  string `json:"address" validate:"omitempty,min=5"`
	Phone    string `json:"phone" validate:"omitempty,phone"`
}
