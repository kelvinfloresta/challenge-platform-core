package campaign_case_test

import (
	"conformity-core/fixtures"
	"conformity-core/gateways/challenge_gateway"
	"conformity-core/gateways/department_gateway"
	"conformity-core/usecases/campaign_case"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Create__Should_fail_if_ChallengeID_is_duplicated(t *testing.T) {
	sut, _ := fixtures.CampaignCase.NewUnit()
	const SAME_ID = "same_challenge_id"
	const COMPANY_ID = "any_company_id"

	_, err := sut.Create(fixtures.FakeBackofficeCtx, campaign_case.CreateCampaignCaseInput{
		Challenges: []campaign_case.CreateScheduledChallenge{
			{
				ChallengeID: SAME_ID,
			},
			{
				ChallengeID: SAME_ID,
			},
		},
		CompanyID: COMPANY_ID,
	})
	require.ErrorIs(t, err, campaign_case.ErrCampaignWithDuplicatedChallengesIds)
}

func Test_Create__Should_not_notify_if_does_not_have_scheduled_challenge_today(t *testing.T) {
	departmentCaseMock, departmentGatewayMock := fixtures.DepartmentCase.NewUnit()
	companyCase, _ := fixtures.CompanyCase.NewUnit()
	challengeCase, challengeGatewayMock := fixtures.ChallengeCase.NewUnit()
	notificationCase, sendEmailMock := fixtures.NotificationCase.NewUnit()
	gatewayMock := &fixtures.CampaignGatewayMock{}

	sut := campaign_case.New(
		gatewayMock,
		departmentCaseMock,
		companyCase,
		challengeCase,
		notificationCase,
	)

	departmentGatewayMock.GetUsersByFilter__Output = &[]department_gateway.GetUsersByFilterGatewayOutput{{
		UserID:       "any_user_id",
		DepartmentID: "any_department_id",
		Name:         "any_name",
		Email:        "any_email",
	}}

	challengeGatewayMock.GetQuestions__Output = &[]challenge_gateway.GetQuestionsOutput{{
		ChallengeID: "any_challenge_id",
		QuestionID:  "any_question_id",
	}}

	tomorrow := time.Now().AddDate(0, 0, 1)
	_, err := sut.Create(fixtures.FakeBackofficeCtx, campaign_case.CreateCampaignCaseInput{
		Title:     "Campaign Test",
		CompanyID: "any_company_id",
		Challenges: []campaign_case.CreateScheduledChallenge{{
			ChallengeID: "any_challenge_id",
			StartDate:   tomorrow,
			EndDate:     tomorrow.AddDate(0, 0, 7),
		}},
	})

	require.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
	require.Len(t, sendEmailMock.Messages, 0)
}
