package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ChallengeOptionJSON struct {
	ID      string
	Title   string
	Correct bool
}

type ChallengeQuestionJSON struct {
	ID      string
	Title   string
	Options []ChallengeOptionJSON
}

type ChallengeMediaJSON struct {
	ID          string
	Title       string
	Path        string
	Description string
}

type Challenge struct {
	ID        string `gorm:"type:uuid"`
	Title     string `gorm:"not null"`
	Segment   string `gorm:"not null; default:''"`
	Media     datatypes.JSON
	Questions datatypes.JSON
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (c Challenge) TableName() string {
	return "challenges"
}
