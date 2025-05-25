package models

import "time"

type RefreshToken struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    Token     string    `gorm:"not null;unique"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}