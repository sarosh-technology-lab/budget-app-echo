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

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (productService ProductService) List(products []*models.Product, pagination *common.Pagination) (*common.Pagination, error) {
	productService.DB.Scopes(pagination.Paginate()).Joins("Category").Joins("SubCategory").Find(&products)
	pagination.Items = products

	return pagination, nil
}

func (productService ProductService) GetById (id uint) (*models.Product, error) {
	 var product *models.Product
	 result := productService.DB.First(&product, id)
	 if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, custom_app_errors.NewNotFoundError("product not found")
		}
		return nil, errors.New("failed to fetch product")
	 }

	 return product, nil
}

func (productService ProductService) Create(data requests.ProductRequestable) (*models.Product, error) {
	name := data.GetName()
	categoryId := data.GetCategoryId()
    subCategoryId := data.GetSubCategoryId()
    description := data.GetDescription()
    image := data.GetImage()
    
	slug := strings.ToLower(name)
	slug = strings.Replace(slug, " ", "_", -1)
	product := &models.Product{
		CategoryId: categoryId,
        SubCategoryId: subCategoryId,
		Slug: slug,
		Name: name,
        Description: description,
        Image: image,
	}

	 _, err := NewCategoryService(productService.DB).GetById(categoryId) // check if category exists
	if err != nil {
		 return nil, errors.New("category does not exist")
	}

	// Start a transaction
	err = productService.DB.Transaction(func(tx *gorm.DB) error {
		// Use the transaction (tx) instead of the main DB connection
		result := tx.Where(models.Product{Slug: slug, Name: name, CategoryId: categoryId, SubCategoryId: subCategoryId, Description: description, Image: image}).FirstOrCreate(product)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				return errors.New("product already exists")
			}
			return errors.New("failed to create product")
		}
		// If no error, transaction commits automatically
		return nil
	})

	// Handle the transaction result
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s ProductService) Update(data requests.ProductRequestable) error {
    id := data.GetId()
    if id == 0 {
        return errors.New("invalid product ID")
    }

    name := data.GetName()
    categoryId := data.GetCategoryId()

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
    var existing models.Product
    if err := s.DB.First(&existing, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("product does not exist")
        }
        return fmt.Errorf("failed to find product: %w", err)
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

        if len(updates) == 0 {
            return nil // nothing to do
        }

        // Check for uniqueness conflict (excluding self)
        var conflict models.Product
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
            return errors.New("another product with the same name/slug/category already exists")
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }

        // Now perform the safe update by ID
        result := tx.Model(&models.Product{}).
            Where("id = ?", id).
            Updates(updates)

        if result.Error != nil {
            if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
                return errors.New("update would violate unique constraint")
            }
            return fmt.Errorf("failed to update product: %w", result.Error)
        }

        if result.RowsAffected == 0 {
            return errors.New("product not found or no changes applied")
        }

        return nil
    })
}


func (productService ProductService) DeleteById (id uint) (error) {
	 var product *models.Product
	 product, err := productService.GetById(id)
	if err != nil {
		 return err
	}

	productService.DB.Delete(product)

	 return nil
}