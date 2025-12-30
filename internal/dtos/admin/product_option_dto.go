package admin

type ProductOptionDTO struct {
	ProductID                     uint    `json:"product_id" validate:"required,gt=0"`
	OptionType                    string  `json:"option_type" validate:"required,oneof=color size material extra"`
	OptionName                    string  `json:"option_name" validate:"required,min=1"`
	PriceModifierType             string  `json:"price_modifier_type" validate:"required,oneof=absolute percent"`
	PriceModifierValue            float64 `json:"price_modifier_value" validate:"required"`
	ProductionTimeModifierDays    int     `json:"production_time_modifier_days" validate:"omitempty"`
	ProductionTimeModifierPercent *int    `json:"production_time_modifier_percent" validate:"omitempty"`
}
