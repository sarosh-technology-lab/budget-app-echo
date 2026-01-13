package services

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/common"
	"budget-backend/internal/custom_app_errors"
	"budget-backend/internal/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SubCategoryService struct {
	DB *gorm.DB
}

func NewSubCategoryService(db *gorm.DB) *SubCategoryService {
	return &SubCategoryService{DB: db}
}

func (subCategoryService SubCategoryService) List(subCategories []*models.SubCategory, pagination *common.Pagination) (*common.Pagination, error) {
	subCategoryService.DB.Scopes(pagination.Paginate()).Preload("Category").Find(&subCategories)
	pagination.Items = subCategories

	return pagination, nil
}

func (subCategoryService SubCategoryService) GetById (id uint) (*models.SubCategory, error) {
	 var subCategory *models.SubCategory
	 result := subCategoryService.DB.First(&subCategory, id)
	 if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_app_errors.NewNotFoundError("sub category not found")
		}
		return nil, errors.New("failed to fetch sub category")
	 }

	 return subCategory, nil
}

func (subCategoryService SubCategoryService) Create(data requests.SubCategoryRequestable) (*models.SubCategory, error) {
	name := data.GetName()
	categoryId := data.GetCategoryId()
	isCustom := data.GetIsCustom()
	slug := strings.ToLower(name)
	slug = strings.Replace(slug, " ", "_", -1)
	subCategory := &models.SubCategory{
		CategoryId: categoryId,
		Slug: slug,
		Name: name,
		IsCustom: isCustom,
	}

	 _, err := NewCategoryService(subCategoryService.DB).GetById(categoryId) // check if category exists
	if err != nil {
		 return nil, errors.New("category does not exist")
	}

	// Start a transaction
	err = subCategoryService.DB.Transaction(func(tx *gorm.DB) error {
		// Use the transaction (tx) instead of the main DB connection
		result := tx.Where(models.SubCategory{Slug: slug, Name: name, CategoryId: categoryId}).FirstOrCreate(subCategory)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				return errors.New("sub category already exists")
			}
			return errors.New("failed to create sub category")
		}
		// If no error, transaction commits automatically
		return nil
	})

	// Handle the transaction result
	if err != nil {
		return nil, err
	}

	return subCategory, nil
}

func (s SubCategoryService) Update(data requests.SubCategoryRequestable) error {
    id := data.GetId()
    if id == 0 {
        return errors.New("invalid sub-category ID")
    }

    name := data.GetName()
    categoryId := data.GetCategoryId()
    // isCustom := data.GetIsCustom()

    // Optional: skip if nothing meaningful changed (depends on your request struct)
    if name == "" && categoryId == 0 {
        return nil // or return error — your choice
    }

    slug := ""
    if name != "" {
        slug = strings.ToLower(name)
        slug = strings.ReplaceAll(slug, " ", "_") // cleaner than Replace with -1
    }

    // Fetch existing to verify it exists (and to get old values if needed)
    var existing models.SubCategory
    if err := s.DB.First(&existing, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("sub category does not exist")
        }
        return fmt.Errorf("failed to find sub category: %w", err)
    }

    // Validate new category if being changed
    if categoryId != 0 && categoryId != existing.CategoryId {
        if _, err := NewCategoryService(s.DB).GetById(categoryId); err != nil {
            return errors.New("category does not exist")
        }
    }

    return s.DB.Transaction(func(tx *gorm.DB) error {
        // Prepare updates (only include fields that are actually changing / provided)
        updates := map[string]interface{}{}

        if name != "" {
            updates["Name"] = name
            updates["Slug"] = slug
        }
        if categoryId != 0 {
            updates["CategoryId"] = categoryId
        }
        // if isCustom != nil { // assuming GetIsCustom returns *bool
        //     updates["IsCustom"] = *isCustom
        // }

        if len(updates) == 0 {
            return nil // nothing to do
        }

        // Check for uniqueness conflict (excluding self)
        var conflict models.SubCategory
        q := tx.Where("id != ?", id)

        if name != "" {
            q = q.Where("slug = ? AND name = ?", slug, name)
        } else {
            // If name not changing, check against existing name/slug + new category if changed
            q = q.Where("slug = ? AND name = ?", existing.Slug, existing.Name)
        }

        if categoryId != 0 {
            q = q.Where("category_id = ?", categoryId)
        } else {
            q = q.Where("category_id = ?", existing.CategoryId)
        }

        if err := q.First(&conflict).Error; err == nil {
            return errors.New("another sub-category with the same name/slug/category already exists")
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }

        // Now perform the safe update by ID
        result := tx.Model(&models.SubCategory{}).
            Where("id = ?", id).
            Updates(updates)

        if result.Error != nil {
            if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
                return errors.New("update would violate unique constraint")
            }
            return fmt.Errorf("failed to update sub-category: %w", result.Error)
        }

        if result.RowsAffected == 0 {
            return errors.New("sub-category not found or no changes applied")
        }

        return nil
    })
}


func (subCategoryService SubCategoryService) DeleteById (id uint) (error) {
	 var subCategory *models.SubCategory
	 subCategory, err := subCategoryService.GetById(id)
	if err != nil {
		 return err
	}

	subCategoryService.DB.Delete(subCategory)

	 return nil
}