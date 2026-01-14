package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CategoryId uint    `gorm:"index;not null" json:"category_id"`
	SubCategoryId uint    `gorm:"index;not null" json:"sub_category_id"`
	Name string    `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Slug string  `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Image string `gorm:"type:varchar(255);not null" json:"image"`
	CreatedAt time.Time `json:"created_at" json:"created_at"` // GORM auto-manages
	UpdatedAt time.Time `json:"updated_at" json:"updated_at"` // GORM auto-manages
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // GORM auto-manages (soft delete)
	Category    Category  `gorm:"foreignKey:CategoryId;references:ID" json:"category"`
	SubCategory    SubCategory  `gorm:"foreignKey:SubCategoryId;references:ID" json:"sub_category"`
}