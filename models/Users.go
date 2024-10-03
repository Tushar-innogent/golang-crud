package models

import (
	"time"

	"gorm.io/gorm"
)

// User has many posts
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Name      string `gorm:"size:100;not null"`
	Password  string `json:"password" validate:"required"`
	Role      Role   `json:"role"`
	Email     string `gorm:"size:100;unique;not null" validate:"required,email"`
	CompanyID uint
	Company   Company
	Posts     []Post `gorm:"constraint:OnDelete:CASCADE;"`
}

// BeforeSave callback to set the role as a string when saving to the DB
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// Convert Role enum to string when saving
	u.Role = ParseRole(u.Role.String())
	return
}
