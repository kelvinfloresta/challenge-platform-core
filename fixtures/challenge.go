package fixtures

import (
	"conformity-core/gateways/challenge_gateway"
	"conformity-core/usecases/challenge_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type challengeCaseFixture struct{}

var ChallengeCase challengeCaseFixture = challengeCaseFixture{}

func (u challengeCaseFixture) NewUnit() (*challenge_case.ChallengeCase, *ChallengeGatewayMock) {
	gatewayMock := &ChallengeGatewayMock{}
	challengeCase := challenge_case.New(gatewayMock)

	return challengeCase, gatewayMock
}

func (u challengeCaseFixture) NewIntegration() *challenge_case.ChallengeCase {
	gateway := &challenge_gateway.GormChallengeGatewayFacade{DB: DB_Test}
	return challenge_case.New(gateway)
}

type ChallengeGatewayMock struct {
	Create__Input        challenge_gateway.CreateChallengeGatewayInput
	GetById__Output      *challenge_gateway.GetByIdChallengeGatewayOutput
	GetOption__Output    *challenge_gateway.Option
	GetQuestions__Output *[]challenge_gateway.GetQuestionsOutput
}

func (c *ChallengeGatewayMock) List() ([]challenge_gateway.ListOutput, error) {
	return nil, nil
}

func (u *ChallengeGatewayMock) Create(input challenge_gateway.CreateChallengeGatewayInput) (string, error) {
	u.Create__Input = input
	return "id", nil
}

func (u *ChallengeGatewayMock) GetById(input string) (*challenge_gateway.GetByIdChallengeGatewayOutput, error) {
	return u.GetById__Output, nil
}

func (u *ChallengeGatewayMock) GetOption(input challenge_gateway.GetOptionInput) (*challenge_gateway.Option, error) {
	return u.GetOption__Output, nil
}

func (u *ChallengeGatewayMock) GetQuestions(challengesIds []string) (*[]challenge_gateway.GetQuestionsOutput, error) {
	return u.GetQuestions__Output, nil
}

func CreateChallenge(t *testing.T) *challenge_gateway.GetByIdChallengeGatewayOutput {
	challengeCase := ChallengeCase.NewIntegration()

	id, err := challengeCase.Create(DUMMY_CONTEXT, challenge_case.CreateChallengeCaseInput{
		Title: "Challenge 1",
		Media: challenge_gateway.CreateMedia{
			Title: "Media 1", Path: "Path", Description: "Desc",
		},
		Questions: []challenge_gateway.CreateQuestion{
			{Title: "Question 1", Options: []challenge_gateway.CreateOption{
				{Title: "Option 1", Correct: true},
				{Title: "Option 2", Correct: false},
				{Title: "Option 3", Correct: false},
				{Title: "Option 4", Correct: false},
			}},
		},
	})
	require.Nil(t, err)
	challenge, err := challengeCase.GetById(id)
	require.Nil(t, err)

	return challenge
}
