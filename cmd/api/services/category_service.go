package services

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/internal/custom_app_errors"
	"budget-backend/internal/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

func (categoryService CategoryService) List() ([]*models.Category, error) {
	var categories []*models.Category
	result := categoryService.DB.Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}

	return categories, nil
}

func (categoryService CategoryService) GetById (id uint) (*models.Category, error) {
	 var category *models.Category
	 result := categoryService.DB.First(&category, id)
	 if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_app_errors.NewNotFoundError("category not found")
		}
		return nil, errors.New("failed to fetch category")
	 }

	 return category, nil
}

func (categoryService CategoryService) Create(data requests.CategoryRequest) (*models.Category, error) {
	slug := strings.ToLower(data.Name)
	slug = strings.Replace(slug, " ", "_", -1)
	category := &models.Category{
		Slug: slug,
		Name: data.Name,
		IsCustom: data.IsCustom,
	}

	result := categoryService.DB.Where(models.Category{Slug: slug, Name: data.Name}).FirstOrCreate(category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return category, nil
	}
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

func (categoryService CategoryService) DeleteById (id uint) (error) {
	 var category *models.Category
	 category, err := categoryService.GetById(id)
	if err != nil {
		 return err
	}

	categoryService.DB.Delete(category)

	 return nil
}