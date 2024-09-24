package models

import "time"

// User has many posts
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;unique;not null" validate:"required,email"`
	CompanyID uint
	Company   Company	
	Posts     []Post `gorm:"constraint:OnDelete:CASCADE;"`
}
