package fixtures

import (
	"conformity-core/gateways/challenge_gateway"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DefaultCampaign struct {
	CompanyID          string
	DepartmentID       string
	UserID             string
	CampaignID         string
	Challenge          *challenge_gateway.GetByIdChallengeGatewayOutput
	StartDate          time.Time
	EndDate            time.Time
	ScheduledChallenge []campaign_case.CreateScheduledChallenge
}

func GetCorrectOption(question challenge_gateway.Question) challenge_gateway.Option {
	return question.Options[0]
}

func GetWrongOption(question challenge_gateway.Question) challenge_gateway.Option {
	return question.Options[1]
}

func CreateDefaultCampaign(t *testing.T, startDate *time.Time) *DefaultCampaign {
	if startDate == nil {
		t := time.Now()
		startDate = &t
	}

	campaignCase := CampaignCase.NewIntegration()

	user := CreateUser(t, nil)

	challenge := CreateChallenge(t)

	endDate := utils.EndOfDay(startDate.AddDate(0, 0, 7))

	challengesToSchedule := []campaign_case.CreateScheduledChallenge{
		{
			ChallengeID: challenge.ID,
			StartDate:   *startDate,
			EndDate:     endDate,
		},
	}

	campaignId, err := campaignCase.Create(FakeBackofficeCtx, campaign_case.CreateCampaignCaseInput{
		Title:      "Campaign 1",
		Challenges: challengesToSchedule,
		CompanyID:  user.CompanyID,
	})

	assert.Nil(t, err)

	return &DefaultCampaign{
		CompanyID:          user.CompanyID,
		DepartmentID:       user.DepartmentID,
		UserID:             user.ID,
		CampaignID:         campaignId,
		Challenge:          challenge,
		StartDate:          *startDate,
		EndDate:            endDate,
		ScheduledChallenge: challengesToSchedule,
	}
}
