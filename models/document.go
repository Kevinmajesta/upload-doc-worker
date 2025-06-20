package models

import (
	"time"
)

type Document struct {
    Id         uint      `gorm:"primaryKey" json:"id"`
    Title      string    `json:"title"`
    FilePath   string    `json:"file_path"`
    EmployeeID uint      `json:"employee_id"`
    Employee   Employee  `json:"employee"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}



