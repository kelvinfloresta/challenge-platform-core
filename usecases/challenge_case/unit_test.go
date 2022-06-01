package challenge_case_test

import (
	"conformity-core/fixtures"
	"conformity-core/gateways/challenge_gateway"
	"conformity-core/usecases/challenge_case"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Answer__Should_fail_if_option_not_found(t *testing.T) {
	sut, gateway := fixtures.ChallengeCase.NewUnit()
	gateway.GetOption__Output = nil

	_, err := sut.Answer(fixtures.DUMMY_CONTEXT, challenge_case.AnswerInput{
		ChallengeID: "any_challenge_id",
		OptionID:    "any_option_id",
	})

	assert.ErrorIs(t, err, challenge_case.ErrOptionNotFound)
}

func Test_Answer__Should_return_true_if_option_is_correct(t *testing.T) {
	sut, gateway := fixtures.ChallengeCase.NewUnit()
	gateway.GetOption__Output = &challenge_gateway.Option{
		ID:      "any_id",
		Title:   "Golang is better",
		Correct: true,
	}

	result, err := sut.Answer(fixtures.DUMMY_CONTEXT, challenge_case.AnswerInput{
		ChallengeID: "any_challenge_id",
		OptionID:    "any_option_id",
	})

	assert.Nil(t, err)
	assert.True(t, result)
}

func Test_Answer__Should_return_false_if_option_is_not_correct(t *testing.T) {
	sut, gateway := fixtures.ChallengeCase.NewUnit()
	gateway.GetOption__Output = &challenge_gateway.Option{
		ID:      "any_id",
		Title:   "Not correct",
		Correct: false,
	}

	result, err := sut.Answer(fixtures.DUMMY_CONTEXT, challenge_case.AnswerInput{
		ChallengeID: "any_challenge_id",
		OptionID:    "any_option_id",
	})

	assert.Nil(t, err)
	assert.False(t, result)
}
