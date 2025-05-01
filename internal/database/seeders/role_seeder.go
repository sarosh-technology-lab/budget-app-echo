package seeders

import (
	"budget-backend/common"
	"budget-backend/internal/models"
)

func SeedRoles() {
	db, err := common.Sql()
	if err != nil {
		panic(err)
	}

	var roles = []models.Role{{Name: "Admin"}, {Name: "User"}}
	db.Create(&roles)
}