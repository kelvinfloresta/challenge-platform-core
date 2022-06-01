package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID                       string `gorm:"type:uuid"`
	Name                     string `gorm:"not null"`
	Document                 string `gorm:"not null"`
	Workspace                string `gorm:"uniqueIndex:unique_workspace"`
	Domain                   string `gorm:"uniqueIndex:unique_domain"`
	IdentityProviderMetadata string
	RequirePassword          bool `gorm:"not null; default: false"`
	SSOEnabled               bool `gorm:"not null; default: false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (c Company) TableName() string {
	return "companies"
}
