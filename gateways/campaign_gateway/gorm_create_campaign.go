package campaign_gateway

import (
	"conformity-core/frameworks/database/gorm/models"

	"gorm.io/gorm"
)

func (g GormCampaignGatewayFacade) Create(input CreateCampaignGatewayInput) (string, error) {
	campaign := &models.Campaign{
		ID:    input.CampaignID,
		Title: input.Title,
	}

	results := adaptCreateResult(input.Results)
	challenges := adaptCreateScheduledChallenge(input.ScheduledChallenges)

	err := g.DB.Conn.Transaction(createCampaign(campaign, results, challenges))

	return input.CampaignID, err
}

func createCampaign(
	campaign *models.Campaign,
	results []*models.CampaignResult,
	challenges []*models.ScheduledChallenge,
) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		err := tx.Create(campaign).Error

		if err != nil {
			return err
		}

		err = tx.Create(results).Error
		if err != nil {
			return err
		}

		return tx.Create(challenges).Error
	}
}
