package database

import (
	"conformity-core/frameworks/database/gorm/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) (err error) {
	return db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Department{},
		&models.UserCompanyDepartment{},
		&models.Challenge{},
		&models.Campaign{},
		&models.CampaignResult{},
		&models.ScheduledChallenge{},
	)
}
