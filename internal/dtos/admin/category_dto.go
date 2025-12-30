package admin

type CategoryDTO struct {
	DepartmentID uint   `json:"department_id" validate:"required,gt=0"`
	Name         string `json:"name" validate:"required,min=2"`
	Description  string `json:"description" validate:"omitempty,min=2"`
}
