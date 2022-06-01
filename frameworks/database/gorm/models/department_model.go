package models

import (
	"conformity-core/enums/user_company_enum"
	"time"

	"gorm.io/gorm"
)

type UserCompanyDepartment struct {
	Login  string                   `gorm:"not null; uniqueIndex:unique_login"`
	Status user_company_enum.Status `gorm:"not null; default:active"`
	Role   user_company_enum.Role   `gorm:"not null; default:user"`
	OID    string

	User   User
	UserID string `gorm:"not null"`

	Department   Department
	DepartmentID string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (u UserCompanyDepartment) TableName() string {
	return "users_companies_departments"
}

type Department struct {
	ID        string `gorm:"type:uuid"`
	Name      string `gorm:"not null"`
	CompanyID string `gorm:"not null"`
	Company   Company
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (u Department) TableName() string {
	return "departments"
}
