package services

import (
    "furniture-shop/internal/models"
)

// Calculate unit price with selected option price modifiers
func CalculateUnitPrice(product models.Product, selected []SelectedOption) float64 {
    price := product.BasePrice
    byID := map[uint]models.ProductOption{}
    for _, o := range product.Options { byID[o.ID] = o }
    for _, so := range selected {
        if opt, ok := byID[so.ID]; ok {
            switch opt.PriceModifierType {
                case "absolute":
                    price += opt.PriceModifierValue
                case "percent":
                    price = price * (1.0 + opt.PriceModifierValue/100.0)
            }
        }
    }
    return price
}

// CalculateUnitPriceWithOptions allows passing options explicitly when product.Options are not populated.
func CalculateUnitPriceWithOptions(base models.Product, options []models.ProductOption, selected []SelectedOption) float64 {
    base.Options = options
    return CalculateUnitPrice(base, selected)
}

