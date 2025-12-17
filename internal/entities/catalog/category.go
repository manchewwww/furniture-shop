package catalog

import "time"

type Category struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    DepartmentID uint      `json:"department_id"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

