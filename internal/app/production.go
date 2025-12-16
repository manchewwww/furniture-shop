package app

import (
    "encoding/json"
    "math"
    "strings"

    "furniture-shop/internal/database"
    models "furniture-shop/internal/domain/entity"
)

type SelectedOption struct {
    ID   uint   `json:"id"`
    Type string `json:"type"`
}

func CalculateItemProductionTime(product models.Product, selected []SelectedOption) int {
    days := product.BaseProductionTimeDays
    // Load options if not already
    if len(product.Options) == 0 {
        var opts []models.ProductOption
        database.DB.Where("product_id = ?", product.ID).Find(&opts)
        product.Options = opts
    }
    optionByID := map[uint]models.ProductOption{}
    for _, o := range product.Options { optionByID[o.ID] = o }

    for _, so := range selected {
        if opt, ok := optionByID[so.ID]; ok {
            days += opt.ProductionTimeModifierDays
            if opt.ProductionTimeModifierPercent != nil {
                days = int(math.Round(float64(days) * (1.0 + float64(*opt.ProductionTimeModifierPercent)/100.0)))
            }
        }
    }
    name := strings.ToLower(product.Name)
    if strings.Contains(name, "180") || strings.Contains(name, "голям") { days += 2 }
    if days < 1 { days = 1 }
    return days
}

func CalculateOrderProductionTime(items []models.OrderItem) int {
    max := 0
    for _, it := range items { if it.CalculatedProductionTimeDays > max { max = it.CalculatedProductionTimeDays } }
    // workload: count orders "впроизводство"
    var cnt int64
    database.DB.Model(&models.Order{}).Where("status = ?", "впроизводство").Count(&cnt)
    if cnt >= 10 { max += 3 } else if cnt >= 5 { max += 1 }
    return max
}

func MarshalSelectedOptions(selected []SelectedOption) string {
    b, _ := json.Marshal(selected)
    return string(b)
}


// CalculateOrderProductionTimeWithWorkload is a DB-free variant used by services layer
func CalculateOrderProductionTimeWithWorkload(items []models.OrderItem, workloadCount int64) int {
    max := 0
    for _, it := range items { if it.CalculatedProductionTimeDays > max { max = it.CalculatedProductionTimeDays } }
    if workloadCount >= 10 { max += 3 } else if workloadCount >= 5 { max += 1 }
    return max
}


