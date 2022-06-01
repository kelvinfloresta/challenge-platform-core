package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	g "conformity-core/gateways/campaign_gateway"
	"encoding/json"
	"time"

	"github.com/getsentry/sentry-go"
)

type AddUserToCampaignsInput struct {
	UserID       string
	DepartmentID string
	CompanyID    string
}

func (c CampaignCase) AddUserToCampaigns(ctx *context.CoreCtx, input AddUserToCampaignsInput) (err error) {
	if input.CompanyID == "" {
		department, err := c.departmentCase.GetById(input.DepartmentID)
		if err != nil {
			return err
		}

		if department == nil {
			sentry.CaptureMessage("IncludeUser Department not found:" + input.DepartmentID)
			return nil
		}
		input.CompanyID = department.CompanyID
	}

	challenges, err := c.gateway.ListChallenges(ctx, g.ListChallengesInput{
		EndDateGTE: time.Now(),
		CompanyID:  input.CompanyID,
	})

	if err != nil {
		return
	}

	if len(challenges) == 0 {
		return
	}

	results := []g.CreateResultsInput{}
	for _, c := range challenges {
		questions := make([]struct{ ID string }, 0)
		if err = json.Unmarshal(c.Questions, &questions); err != nil {
			return
		}

		for _, q := range questions {
			results = append(results, g.CreateResultsInput{
				UserID:       input.UserID,
				DepartmentID: input.DepartmentID,
				ChallengeID:  c.ID,
				CampaignID:   c.CampaignID,
				QuestionID:   q.ID,
				Status:       campaign_enum.Active,
				Tries:        0,
				Correct:      false,
			})
		}
	}

	return c.gateway.CreateResults(results, nil)
}
