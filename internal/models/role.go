package models

import (
    "time"
    "gorm.io/gorm"
)

type Role struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `gorm:"unique" json:"name"`
    CreatedAt time.Time `gorm:"index" json:"created_at"`
    UpdatedAt time.Time `gorm:"index" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}