package admin

type ProductDTO struct {
	CategoryID             uint    `json:"category_id" validate:"required,gt=0"`
	Name                   string  `json:"name" validate:"required,min=2"`
	ShortDescription       string  `json:"short_description" validate:"omitempty,min=2"`
	LongDescription        string  `json:"long_description" validate:"omitempty,min=2"`
	BasePrice              float64 `json:"base_price" validate:"required,gte=0"`
	BaseProductionTimeDays int     `json:"base_production_time_days" validate:"required,gte=0"`
	ImageURL               string  `json:"image_url" validate:"omitempty,url"`
	DefaultWidth           int     `json:"default_width" validate:"required,gt=0"`
	DefaultHeight          int     `json:"default_height" validate:"required,gt=0"`
	DefaultDepth           int     `json:"default_depth" validate:"required,gt=0"`
}
