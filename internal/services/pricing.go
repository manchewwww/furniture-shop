package services

import (
    "furniture-shop/internal/database"
    "furniture-shop/internal/models"
)

// Calculate unit price with selected option price modifiers
func CalculateUnitPrice(product models.Product, selected []SelectedOption) float64 {
    price := product.BasePrice
    if len(product.Options) == 0 {
        var opts []models.ProductOption
        database.DB.Where("product_id = ?", product.ID).Find(&opts)
        product.Options = opts
    }
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

