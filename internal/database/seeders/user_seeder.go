package seeders

import (
	"budget-backend/common"
	"budget-backend/internal/models"
)

func SeedUsers() {
	db, err := common.Sql()
	if err != nil {
		panic(err)
	}

	hashedPassword, err := common.HashPassword("password")
	if err != nil {
		panic(err)
	}

	var users = []models.User{{RoleId: 1, FirstName: "Admin", LastName: "Admin", Email: "admin@example.com", Password: hashedPassword}}
	db.Create(&users)
}
