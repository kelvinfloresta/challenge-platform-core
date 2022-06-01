package campaign_case

import (
	"conformity-core/context"
	"conformity-core/enums/campaign_enum"
	"conformity-core/enums/user_company_enum"
	g "conformity-core/gateways/campaign_gateway"
	"time"
)

type ListUsersInput struct {
	StartDateLTE time.Time
	StartDate    time.Time
	EndDateLTE   time.Time
	EndDateGT    time.Time
	CampaignID   string
	ChallengeID  string
	Correct      *bool
	TriesLT      *uint8
}

type ListUsersOutput struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	DepartmentName string `json:"departmentName"`
	CompanyID      string `json:"-"`
	UserID         string `json:"userId"`
	DepartmentID   string `json:"departmentId"`
	IMPD           uint8  `json:"impd"`
}

func (u CampaignCase) ListUsers(ctx *context.CoreCtx, input ListUsersInput) ([]ListUsersOutput, error) {
	userCompanyStatus := user_company_enum.Active
	campaignResultStatus := campaign_enum.Active

	filter := g.ListUsersInput{
		StartDateLTE:         input.StartDateLTE,
		StartDate:            input.StartDate,
		EndDateLTE:           input.EndDateLTE,
		EndDateGT:            input.EndDateGT,
		TriesLessThan:        input.TriesLT,
		Correct:              input.Correct,
		UseCompanyStatus:     &userCompanyStatus,
		CampaignResultStatus: &campaignResultStatus,
		ChallengeID:          input.ChallengeID,
		CampaignID:           input.CampaignID,
	}

	data, err := u.gateway.ListUsers(ctx, filter)

	if err != nil {
		return nil, err
	}

	result := []ListUsersOutput{}

	for _, lqo := range data {
		result = append(result, ListUsersOutput{
			Email:          lqo.Email,
			Name:           lqo.Name,
			DepartmentName: lqo.DepartmentName,
			CompanyID:      lqo.CompanyID,
			UserID:         lqo.UserID,
			DepartmentID:   lqo.DepartmentID,
			IMPD:           lqo.IMPD,
		})
	}

	return result, nil

}
