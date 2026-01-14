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

     _, err = NewSubCategoryService(productService.DB).GetById(subCategoryId) // check if sub category exists
	if err != nil {
		 return nil, errors.New("sub category does not exist")
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
    subCategoryId := data.GetSubCategoryId()
    description := data.GetDescription()
    image := data.GetImage()

    if name == "" && categoryId == 0 && subCategoryId == 0 && description == "" && image == "" {
        return nil
    }

    slug := ""
    if name != "" {
        slug = strings.ToLower(name)
        slug = strings.ReplaceAll(slug, " ", "_")
    }

    // Fetch existing to verify it exists
    var existing models.Product
    if err := s.DB.First(&existing, int(id)).Error; err != nil {
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

    // Validate new sub category if being changed
    if subCategoryId != 0 && subCategoryId != existing.SubCategoryId {
        if _, err := NewSubCategoryService(s.DB).GetById(subCategoryId); err != nil {
            return errors.New("sub category does not exist")
        }
    }

    return s.DB.Transaction(func(tx *gorm.DB) error {
        // Prepare updates
        updates := map[string]interface{}{}

        if name != "" {
            updates["Name"] = name
            updates["Slug"] = slug
        }
        if categoryId != 0 {
            updates["CategoryId"] = categoryId
        }
        if subCategoryId != 0 {
            updates["SubCategoryId"] = subCategoryId
        }
        if description != "" {
            updates["Description"] = description
        }
        if image != "" {
            updates["Image"] = image
        }

        if len(updates) == 0 {
            return nil
        }

        // Check for uniqueness conflict (excluding self)
        var conflict models.Product
        checkName := name
        checkSlug := slug
        checkCategoryId := categoryId

        if name == "" {
            checkName = existing.Name
            checkSlug = existing.Slug
        }
        if categoryId == 0 {
            checkCategoryId = existing.CategoryId
        }

        result := tx.Where("slug = ? AND name = ? AND category_id = ? AND id != ?", 
            checkSlug, checkName, int(checkCategoryId), int(id)).First(&conflict)

        if result.Error == nil {
            return errors.New("another product with the same name/slug/category already exists")
        } else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return result.Error
        }

        // Perform the update
        result = tx.Model(&models.Product{}).Where("id = ?", int(id)).Updates(updates)

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
