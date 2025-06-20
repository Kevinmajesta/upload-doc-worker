package models

import "time"

type Employee struct {
	Id         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"not null"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Phone      string    `json:"phone"`
	Position   string    `json:"position"`
	Department string    `json:"department"` 
	Documents  []Document `gorm:"foreignKey:EmployeeID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
