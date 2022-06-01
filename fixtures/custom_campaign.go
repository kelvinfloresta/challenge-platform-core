package fixtures

import (
	"conformity-core/usecases/campaign_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type CustomCampaing struct {
	UserIDs     []string
	Challenges  []campaign_case.CreateScheduledChallenge
	DepartmentA string
	DepartmentB string
	CompanyID   string
	CampaignID  string
}

type CreateCustomCampaignInput struct {
	UserAmount int
	Challenges []campaign_case.CreateScheduledChallenge
}

func CreateCustomCampaign(
	t *testing.T,
	input CreateCustomCampaignInput,
) *CustomCampaing {
	campaignCase := CampaignCase.NewIntegration()
	companyId := CreateCompany(t)
	departmentA := CreateDepartment(t, companyId)
	departmentB := CreateDepartment(t, companyId)
	userIds := []string{}

	for i := 0; i < input.UserAmount; i++ {
		if i%2 == 0 {
			userIds = append(userIds, CreateUser(t, &departmentA).ID)
		} else {
			userIds = append(userIds, CreateUser(t, &departmentB).ID)
		}
	}

	campaignId, err := campaignCase.Create(FakeBackofficeCtx, campaign_case.CreateCampaignCaseInput{
		Title:      "Custom Campaign 1",
		Challenges: input.Challenges,
		CompanyID:  companyId,
	})

	require.Nil(t, err)

	return &CustomCampaing{
		UserIDs:     userIds,
		DepartmentA: departmentA,
		DepartmentB: departmentB,
		CompanyID:   companyId,
		CampaignID:  campaignId,
		Challenges:  input.Challenges,
	}
}
