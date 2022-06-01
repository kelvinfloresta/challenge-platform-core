package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string  `gorm:"type:uuid"`
	Name        string  `gorm:"not null; index:idx_name"`
	Password    string  `gorm:"not null"`
	Email       *string `gorm:"uniqueIndex:unique_email"`
	Phone       string  `gorm:"not null"`
	BirthDate   string  `gorm:"not null"`
	Document    string  `gorm:"not null"`
	JobPosition string  `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (u User) TableName() string {
	return "users"
}
