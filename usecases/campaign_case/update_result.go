package campaign_case

import (
	"conformity-core/enums/campaign_enum"
	c_gateway "conformity-core/gateways/campaign_gateway"
	"time"

	"gorm.io/gorm"
)

type UpdateResultInput struct {
	UserID       string
	DepartmentID string
	NewStatus    campaign_enum.ResultStatus
}

func (c CampaignCase) UpdateResult(input UpdateResultInput, tx *gorm.DB) (bool, error) {
	actualCorrect := false
	actualTriesLessThan := uint8(3)
	now := time.Now()

	return c.gateway.UpdateResult(c_gateway.UpdateResultInput{
		NewStatus: input.NewStatus,

		UserID:              input.UserID,
		DepartmentID:        input.DepartmentID,
		ActualCorrect:       &actualCorrect,
		EndDateGreaterThan:  &now,
		ActualTriesLessThan: &actualTriesLessThan,
	},
		tx,
	)
}
