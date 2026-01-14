package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    CategoryId    uint           `gorm:"index;not null" json:"category_id"`
    SubCategoryId uint           `gorm:"index;not null" json:"sub_category_id"`
    
    // Unique index for Name + DeletedAt
    Name          string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_product_name_deleted" json:"name"`
    
    Description   string         `gorm:"type:text;not null" json:"description"`
    
    // Unique index for Slug + DeletedAt
    Slug          string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_product_slug_deleted" json:"slug"`
    
    Image         string         `gorm:"type:varchar(255);not null" json:"image"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    
    // Include DeletedAt in both composite indexes
    DeletedAt     gorm.DeletedAt `gorm:"index;uniqueIndex:idx_product_name_deleted;uniqueIndex:idx_product_slug_deleted" json:"deleted_at"`
    
    Category      Category       `gorm:"foreignKey:CategoryId;references:ID" json:"category"`
    SubCategory   SubCategory    `gorm:"foreignKey:SubCategoryId;references:ID" json:"sub_category"`
}