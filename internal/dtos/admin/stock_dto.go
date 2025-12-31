package admin

type StockDTO struct {
	MaterialName      string  `json:"material_name" validate:"required,min=1"`
	QuantityAvailable float64 `json:"quantity_available" validate:"required,gte=0"`
	Unit              string  `json:"unit" validate:"required,min=1"`
}
