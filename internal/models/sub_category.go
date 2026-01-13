package models

import (
	"time"

	"gorm.io/gorm"
)

type SubCategory struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CategoryId uint    `gorm:"index;not null" json:"category_id"`
	Name string    `gorm:"type:varchar(255);unique;not null" json:"name"`
	Slug string  `gorm:"type:varchar(255);unique;not null" json:"slug"`
	IsCustom  bool    `gorm:"type:bool;not null;default:0" json:"is_custom"`
	CreatedAt time.Time `json:"created_at" json:"created_at"` // GORM auto-manages
	UpdatedAt time.Time `json:"updated_at" json:"updated_at"` // GORM auto-manages
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // GORM auto-manages (soft delete)
	Category    Category  `gorm:"foreignKey:CategoryId;references:ID" json:"category"`
}