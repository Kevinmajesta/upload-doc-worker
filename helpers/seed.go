package helpers

import (
	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/models"
	"time"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	admins := []models.Admin{
		{
			Username: "admin",
			Password: "admin123",
		},
	}

	for i, admin := range admins {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			panic("Failed to hash password: " + err.Error())
		}
		admins[i].Password = string(hashedPassword)

		database.DB.FirstOrCreate(&admins[i], models.Admin{Username: admins[i].Username})
	}
}

func SeedEmployees() {
	employees := []models.Employee{
		{
			Name:       "Kevin Majesta",
			Email:      "kevin@example.com",
			Phone:      "08123456789",
			Position:   "Software Engineer",
			Department: "IT",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Name:       "Anna Sutrisno",
			Email:      "anna@example.com",
			Phone:      "08234567890",
			Position:   "Product Manager",
			Department: "Product",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Name:       "Budi Santoso",
			Email:      "budi@example.com",
			Phone:      "08345678901",
			Position:   "HR Specialist",
			Department: "HR",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, emp := range employees {
		database.DB.FirstOrCreate(&emp, models.Employee{Email: emp.Email})
	}
}
