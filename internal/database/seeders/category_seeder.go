package seeders

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
)

func SeedCategories() {
	db, err := common.Sql()
	if err != nil {
		panic(err)
	}
	categoryService := services.NewCategoryService(db)

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
	}
}