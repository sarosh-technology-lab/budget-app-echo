package services

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/common"
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

func (categoryService CategoryService) List(categories []*models.Category, pagination *common.Pagination) (*common.Pagination, error) {
	categoryService.DB.Scopes(pagination.Paginate()).Find(&categories)
	pagination.Items = categories

	return pagination, nil
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

	// Start a transaction
	err := categoryService.DB.Transaction(func(tx *gorm.DB) error {
		// Use the transaction (tx) instead of the main DB connection
		result := tx.Where(models.Category{Slug: slug, Name: data.Name}).FirstOrCreate(category)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				return errors.New("category already exists")
			}
			return errors.New("failed to create category")
		}
		// If no error, transaction commits automatically
		return nil
	})

	// Handle the transaction result
	if err != nil {
		return nil, err
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