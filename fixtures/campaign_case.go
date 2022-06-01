package fixtures

import (
	coreContext "conformity-core/context"
	g "conformity-core/gateways/campaign_gateway"
	"conformity-core/usecases/campaign_case"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type campaignCaseFixture struct{}

var CampaignCase campaignCaseFixture = campaignCaseFixture{}

func (u campaignCaseFixture) NewIntegration() campaign_case.CampaignCase {
	gateway := &g.GormCampaignGatewayFacade{DB: DB_Test}
	departmentCase := DepartmentCase.NewIntegration()
	companyCase := CompanyCase.NewIntegration()
	challengeCase := ChallengeCase.NewIntegration()
	notificationCase, _ := NotificationCase.NewUnit()
	return campaign_case.New(gateway, departmentCase, companyCase, challengeCase, notificationCase)
}

func (u campaignCaseFixture) NewUnit() (campaign_case.CampaignCase, *CampaignGatewayMock) {
	gatewayMock := &CampaignGatewayMock{}
	departmentCase, _ := DepartmentCase.NewUnit()
	companyCase, _ := CompanyCase.NewUnit()
	challengeCase, _ := ChallengeCase.NewUnit()
	notificationCase, _ := NotificationCase.NewUnit()
	campaignCase := campaign_case.New(gatewayMock, departmentCase, companyCase, challengeCase, notificationCase)
	return campaignCase, gatewayMock
}

type CampaignGatewayMock struct {
	Create__Input                 g.CreateCampaignGatewayInput
	GetById__Output               *g.GetByIdCampaignGatewayOutput
	GetActiveCampaigns__Output    []g.GetActiveCampaignsOutput
	GetChallengesCampaign__Output []g.GetChallengesCampaignOutput
	UpdateResult__Error           error
	List__Output                  []*g.ListOutput
	ListUsers__Output             []g.ListUsersOutput
	ChangeUserStatus__Impl        func(input g.ChangeUserStatusInput, tx *gorm.DB) (bool, error)
	Transaction__Impl             func(fn func(tx *gorm.DB) error) error
	RemoveUser__Impl              func(ctx *coreContext.CoreCtx, input g.RemoveUserInput, tx *gorm.DB) (bool, error)
}

func (u *CampaignGatewayMock) Create(input g.CreateCampaignGatewayInput) (string, error) {
	u.Create__Input = input
	return "id", nil
}

func (u *CampaignGatewayMock) GetById(input string) (*g.GetByIdCampaignGatewayOutput, error) {
	return u.GetById__Output, nil
}

func (c *CampaignGatewayMock) CreateResults(input []g.CreateResultsInput, tx *gorm.DB) error {
	return nil
}

func (c *CampaignGatewayMock) Transaction(fc func(tx *gorm.DB) error) (err error) {
	if c.Transaction__Impl != nil {
		return c.Transaction__Impl(fc)
	}
	return nil
}

func (c *CampaignGatewayMock) UpdateResult(input g.UpdateResultInput, tx *gorm.DB) (bool, error) {
	return false, c.UpdateResult__Error
}

func (c *CampaignGatewayMock) GetResult(ctx *coreContext.CoreCtx, input g.GetResultInput) (*g.Result, error) {
	return nil, nil
}

func (c *CampaignGatewayMock) GetAVGResult(ctx *coreContext.CoreCtx, input g.GetAVGResultInput) ([]*g.GetAVGResultOutput, error) {
	return nil, nil
}

func (u *CampaignGatewayMock) GetActiveCampaigns(ctx *coreContext.CoreCtx) ([]g.GetActiveCampaignsOutput, error) {
	return u.GetActiveCampaigns__Output, nil
}

func (u *CampaignGatewayMock) ListChallenges(ctx *coreContext.CoreCtx, input g.ListChallengesInput) ([]g.GetChallengesCampaignOutput, error) {
	return u.GetChallengesCampaign__Output, nil
}

func (u *CampaignGatewayMock) ListQuestions(ctx *coreContext.CoreCtx, input g.ListQuestionsInput) ([]g.ListQuestionsOutput, error) {
	return nil, nil
}

func (u *CampaignGatewayMock) ListResults(ctx *coreContext.CoreCtx, input g.ListResultsInput) ([]g.ListResultsOutput, error) {
	return []g.ListResultsOutput{}, nil
}

func (u *CampaignGatewayMock) List(ctx *coreContext.CoreCtx, input g.ListInput) ([]*g.ListOutput, error) {
	return u.List__Output, nil
}

func (u *CampaignGatewayMock) ListUsers(ctx *coreContext.CoreCtx, input g.ListUsersInput) ([]g.ListUsersOutput, error) {
	return u.ListUsers__Output, nil
}

func (u *CampaignGatewayMock) ChangeUserStatus(input g.ChangeUserStatusInput, tx *gorm.DB) (bool, error) {
	if u.ChangeUserStatus__Impl != nil {
		return u.ChangeUserStatus__Impl(input, tx)
	}

	return false, nil
}

func (u *CampaignGatewayMock) RemoveUser(ctx *coreContext.CoreCtx, input g.RemoveUserInput, tx *gorm.DB) (bool, error) {
	if u.RemoveUser__Impl != nil {
		return u.RemoveUser__Impl(ctx, input, tx)
	}
	return false, nil
}

func CreateCampaign(t *testing.T) string {
	campaignCase := CampaignCase.NewIntegration()

	id, err := campaignCase.Create(DUMMY_CONTEXT, campaign_case.CreateCampaignCaseInput{
		Title: "Any name",
	})
	require.Nil(t, err)

	return id
}

func AsyncAnswerChallenge(t *testing.T, ctx *coreContext.CoreCtx, campaignCase campaign_case.CampaignCase) (
	func(ctx *coreContext.CoreCtx, input campaign_case.AnswerChallengeInput), *sync.WaitGroup) {

	wg := &sync.WaitGroup{}

	job := func(ctx *coreContext.CoreCtx, input campaign_case.AnswerChallengeInput) {
		// #nosec errcheck
		_, err := campaignCase.AnswerChallenge(ctx, input)
		if err == nil {
			wg.Done()
		} else {
			wg.Done()
		}
	}

	return job, wg
}
