package models

import "time"

type User struct {
	ID uint64
	FirstName string    `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(50);not null" json:"last_name"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;default:null" json:"updated_at"`
}