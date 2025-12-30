package admin

type DepartmentDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description" validate:"omitempty,min=2"`
	ImageURL    string `json:"image_url" validate:"omitempty,url"`
}
