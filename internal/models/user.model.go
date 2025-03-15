package models

import "time"

type User struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(50);not null" json:"last_name"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone     string    `gorm:"type:varchar(11);unique" json:"phone"`
	Address   string    `gorm:"type:text" json:"address"`
	Gender    *string    `gorm:"type:enum('M','F','O');null" json:"gender"` // '*string' represents nullable 
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `json:"created_at"` // GORM auto-manages
	UpdatedAt time.Time `json:"updated_at"` // GORM auto-manages
}