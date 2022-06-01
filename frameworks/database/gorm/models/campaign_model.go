package models

import (
	"conformity-core/enums/campaign_enum"
	"time"

	"gorm.io/gorm"
)

type CampaignResult struct {
	UserID       string `gorm:"not null; uniqueIndex:unique_result"`
	DepartmentID string `gorm:"not null; uniqueIndex:unique_result"`
	ChallengeID  string `gorm:"not null; uniqueIndex:unique_result"`
	CampaignID   string `gorm:"not null; uniqueIndex:unique_result"`
	QuestionID   string `gorm:"not null; uniqueIndex:unique_result"`

	User       User
	Department Department
	Challenge  Challenge

	Tries   uint8                      `gorm:"not null"`
	Correct bool                       `gorm:"not null; default:false"`
	Status  campaign_enum.ResultStatus `gorm:"not null; default:active"`
}

func (c CampaignResult) TableName() string {
	return "campaigns_results"
}

type ScheduledChallenge struct {
	CampaignID  string `gorm:"not null"`
	ChallengeID string `gorm:"not null"`
	Challenge   Challenge

	StartDate time.Time `gorm:"not null; index"`
	EndDate   time.Time `gorm:"not null; index"`
}

func (s ScheduledChallenge) TableName() string {
	return "scheduled_challenges"
}

type Campaign struct {
	ID                  string `gorm:"type:uuid"`
	Title               string `gorm:"not null"`
	Results             []CampaignResult
	ScheduledChallenges []ScheduledChallenge
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

func (c Campaign) TableName() string {
	return "campaigns"
}
