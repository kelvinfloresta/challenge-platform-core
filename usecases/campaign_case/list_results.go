package campaign_case

import (
	"conformity-core/context"
	"conformity-core/domain"
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	"conformity-core/gateways/campaign_gateway"
	"time"
)

type ListResultsInput struct {
	UserID      string `json:"userId"`
	CampaignID  string `json:"campaignId" validate:"required"`
	ChallengeID string `json:"challengeId" validate:"required"`
	CompanyID   string `json:"companyId"`
}

type ListResultsOutput struct {
	UserID         string                   `json:"userId"`
	ChallengeID    string                   `json:"challengeId"`
	CampaignID     string                   `json:"campaignId"`
	DepartmentID   string                   `json:"departmentId"`
	Name           string                   `json:"userName"`
	DepartmentName string                   `json:"departmentName"`
	Tries          uint8                    `json:"tries"`
	Correct        bool                     `json:"correct"`
	ChallengeTitle string                   `json:"challengeTitle"`
	EndDate        time.Time                `json:"endDate"`
	Finished       bool                     `json:"finished"`
	IMPD           uint8                    `json:"impd"`
	Status         user_company_enum.Status `json:"status"`
}

func (u CampaignCase) ListResults(ctx *context.CoreCtx, input ListResultsInput) ([]ListResultsOutput, error) {
	if ctx.Session.Role.IsCommonUser() {
		input.UserID = ctx.Session.UserID
	}

	if !ctx.Session.Role.IsBackoffice() {
		input.CompanyID = ctx.Session.CompanyID
	}

	filter := campaign_gateway.ListResultsInput{
		UserID:      input.UserID,
		CampaignID:  input.CampaignID,
		ChallengeID: input.ChallengeID,
		CompanyID:   input.CompanyID,
		Status:      campaign_enum.Active,
	}

	results, err := u.gateway.ListResults(ctx, filter)

	if err != nil {
		return nil, err
	}

	output := make([]ListResultsOutput, 0, len(results))
	for _, result := range results {
		r := domain.Result(result.Tries, result.Correct)
		output = append(output, ListResultsOutput{
			UserID:         result.UserID,
			ChallengeID:    result.ChallengeID,
			DepartmentID:   result.DepartmentID,
			Name:           result.Name,
			DepartmentName: result.DepartmentName,
			Tries:          result.Tries,
			Correct:        result.Correct,
			ChallengeTitle: result.ChallengeTitle,
			EndDate:        result.EndDate,
			Finished:       r.IsFinished(),
			IMPD:           r.IMPD(),
			Status:         result.Status,
		})
	}

	return output, nil
}
