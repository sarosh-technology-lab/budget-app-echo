package seeders

import (
	"budget-backend/common"
	"budget-backend/internal/models"
)

func SeedPermissions() {
	db, err := common.Sql()
	if err != nil {
		panic(err)
	}

	var permissions = []models.Permission{{Name: "Read:Category"}, {Name: "Create:Category"}, {Name: "Update:Category"}, {Name: "Delete:Category"}, {Name: "Read:User"}, {Name: "Create:User"}, {Name: "Update:User"}, {Name: "Delete:User"}}
	db.Create(&permissions)
}
