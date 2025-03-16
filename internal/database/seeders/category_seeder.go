package main

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
)

func main() {
	db, err := common.Mysql()
	if err != nil {
		panic(err)
	}
	categoryService := services.CategoryService{
		DB: db,
	}

	categories := []string{
		"Food", "Gifts", "Health", "Fashion", "Medical", "Eat Out", "Services", "Information Technology",
	}

	for _, category := range categories {
		_, err = categoryService.Create(requests.CategoryRequest{
			Name: category,
			IsCustom: false,
		})
		if err != nil {
			panic(err.Error())
		}
		println("category " + category + " created")
	}
}