package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(255);unique;not null" json:"name"`
	Slug string  `gorm:"type:varchar(255);unique;not null" json:"slug"`
	IsCustom  bool    `gorm:"type:bool;not null;default:0" json:"is_custom"`
	CreatedAt time.Time `json:"created_at"` // GORM auto-manages
	UpdatedAt time.Time `json:"updated_at"` // GORM auto-manages
	DeletedAt gorm.DeletedAt `gorm:"index"` // GORM auto-manages (soft delete)
}