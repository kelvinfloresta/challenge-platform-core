package campaign_gateway

import (
	"conformity-core/frameworks/database/gorm/models"

	"gorm.io/gorm"
)

func (g GormCampaignGatewayFacade) CreateResults(input []CreateResultsInput, tx *gorm.DB) error {
	if tx == nil {
		tx = g.DB.Conn
	}

	results := make([]models.CampaignResult, 0, len(input))
	for _, r := range input {
		results = append(results, models.CampaignResult{
			UserID:       r.UserID,
			DepartmentID: r.DepartmentID,
			ChallengeID:  r.ChallengeID,
			CampaignID:   r.CampaignID,
			QuestionID:   r.QuestionID,
			Tries:        r.Tries,
			Correct:      r.Correct,
			Status:       r.Status,
		})
	}

	return tx.Create(results).Error
}
