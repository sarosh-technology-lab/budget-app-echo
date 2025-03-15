package models

import "time"

type AppToken struct {
	ID  uint64  `gorm:"primaryKey" json:"id"`
	TargetId   uint `json:"target_id" gorm:"index;not null;"`
	Type string `json:"-" gorm:"index;not null;type:varchar(255)"`
	Token string `json:"-" gorm:"index;not null;type:varchar(255)"`
	Used bool `json:"-" gorm:"index;not null;type:bool"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null;"`
	CreatedAt time.Time `json:"created_at"` // GORM auto-manages
	UpdatedAt time.Time `json:"updated_at"` // GORM auto-manages
}

