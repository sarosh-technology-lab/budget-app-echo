package models

import (
    "time"
    "gorm.io/gorm"
)

type RolePermission struct {
    ID           uint `gorm:"primaryKey" json:"id"`
    RoleID       uint `gorm:"index" json:"role_id"`
    PermissionID uint `gorm:"index" json:"permission_id"`
    CreatedAt time.Time `gorm:"index" json:"created_at"` // GORM auto-manages
    UpdatedAt time.Time `gorm:"index" json:"updated_at"` // GORM auto-manages
    DeletedAt gorm.DeletedAt `gorm:"index"` // GORM auto-manages (soft delete)
}