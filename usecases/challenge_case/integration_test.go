package challenge_case_test

import (
	"conformity-core/fixtures"
	g "conformity-core/gateways/challenge_gateway"
	c_case "conformity-core/usecases/challenge_case"
	"conformity-core/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallengeCase_Create__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.ChallengeCase.NewIntegration()

	mediaToCreate := g.CreateMedia{Title: "Media 1", Path: "path", Description: "desc"}

	questionToCreate := g.CreateQuestion{
		Title: "Question 1",
		Options: []g.CreateOption{
			{Title: "option 1", Correct: true},
			{Title: "option 2", Correct: false},
			{Title: "option 3", Correct: false},
			{Title: "option 4", Correct: false},
		},
	}

	challengeToCreate := c_case.CreateChallengeCaseInput{
		Title:     "Challenge 1",
		Segment:   "Segment",
		Media:     mediaToCreate,
		Questions: []g.CreateQuestion{questionToCreate},
	}

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, challengeToCreate)
	require.Nil(t, err)

	found, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Equal(t, found.ID, id)
	assert.Equal(t, found.Title, challengeToCreate.Title)
	assert.Equal(t, found.Segment, challengeToCreate.Segment)

	mediaFound := found.Media
	assert.Equal(t, mediaToCreate.Title, mediaFound.Title)
	assert.Equal(t, mediaToCreate.Path, mediaFound.Path)
	assert.Equal(t, mediaToCreate.Description, mediaFound.Description)

	questionFound := found.Questions[0]
	assert.Equal(t, questionFound.Title, questionToCreate.Title)

	for i, option := range questionToCreate.Options {
		assert.Equal(t, option.Title, questionToCreate.Options[i].Title)
		assert.Equal(t, option.Correct, questionToCreate.Options[i].Correct)
	}
}

func TestChallengeCase_Create__Should_fail_if_not_have_a_correct_option(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.ChallengeCase.NewIntegration()

	mediaToCreate := g.CreateMedia{Title: "Media 1", Path: "path", Description: "desc"}

	questionToCreate := g.CreateQuestion{
		Title: "Question 1",
		Options: []g.CreateOption{
			{Title: "option 1", Correct: false},
			{Title: "option 2", Correct: false},
			{Title: "option 3", Correct: false},
			{Title: "option 4", Correct: false},
		},
	}

	challengeToCreate := c_case.CreateChallengeCaseInput{
		Title:     "Challenge 1",
		Media:     mediaToCreate,
		Questions: []g.CreateQuestion{questionToCreate},
	}

	_, err := sut.Create(fixtures.DUMMY_CONTEXT, challengeToCreate)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, c_case.ErrWithoutCorrectOption)
}

func TestChallengeCase_Create__Should_fail_if_have_more_than_one_correct_option(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.ChallengeCase.NewIntegration()

	mediaToCreate := g.CreateMedia{Title: "Media 1", Path: "path", Description: "desc"}

	questionToCreate := g.CreateQuestion{
		Title: "Question 1",
		Options: []g.CreateOption{
			{Title: "option 1", Correct: true},
			{Title: "option 2", Correct: true},
			{Title: "option 3", Correct: false},
			{Title: "option 4", Correct: false},
		},
	}

	challengeToCreate := c_case.CreateChallengeCaseInput{
		Title:     "Challenge 1",
		Media:     mediaToCreate,
		Questions: []g.CreateQuestion{questionToCreate},
	}

	_, err := sut.Create(fixtures.DUMMY_CONTEXT, challengeToCreate)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, c_case.ErrMultipleCorrectOptions)
}

func TestChallengeCase_Create__Should_fail_if_questions_is_empty(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.ChallengeCase.NewIntegration()

	mediaToCreate := g.CreateMedia{Title: "Media 1", Path: "path", Description: "desc"}

	challengeToCreate := c_case.CreateChallengeCaseInput{
		Title:     "Challenge 1",
		Media:     mediaToCreate,
		Questions: []g.CreateQuestion{},
	}

	_, err := sut.Create(fixtures.DUMMY_CONTEXT, challengeToCreate)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, c_case.ErrEmptyChallengeQuestions)
}

func TestChallengeCase_Create__Should_fail_if_options_is_empty(t *testing.T) {
	fixtures.CleanTestDatabase()
	EMPTY_OPTIONS := []g.CreateOption{}

	sut := fixtures.ChallengeCase.NewIntegration()

	mediaToCreate := g.CreateMedia{Title: "Media 1", Path: "path", Description: "desc"}
	questionToCreate := g.CreateQuestion{
		Title:   "Question 1",
		Options: EMPTY_OPTIONS,
	}

	challengeToCreate := c_case.CreateChallengeCaseInput{
		Title:     "Challenge 1",
		Media:     mediaToCreate,
		Questions: []g.CreateQuestion{questionToCreate},
	}

	_, err := sut.Create(fixtures.DUMMY_CONTEXT, challengeToCreate)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, c_case.ErrWithoutCorrectOption)
}

func Test_GetOption__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.ChallengeCase.NewIntegration()

	question := challenge.Questions[0]
	correctOption := question.Options[0]

	result, err := sut.Answer(fixtures.DUMMY_CONTEXT, c_case.AnswerInput{
		ChallengeID: challenge.ID,
		QuestionID:  question.ID,
		OptionID:    correctOption.ID,
	})

	assert.Nil(t, err)
	assert.True(t, result)
}

func Test_GetOption__Should_return_false_if_options_is_not_correct(t *testing.T) {
	fixtures.CleanTestDatabase()
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.ChallengeCase.NewIntegration()

	question := challenge.Questions[0]
	wrongOption := question.Options[1]

	result, err := sut.Answer(fixtures.DUMMY_CONTEXT, c_case.AnswerInput{
		ChallengeID: challenge.ID,
		QuestionID:  question.ID,
		OptionID:    wrongOption.ID,
	})

	assert.Nil(t, err)
	assert.False(t, result)
}

func Test_GetOption__Should_return_an_error_if_not_found(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.ChallengeCase.NewIntegration()

	NOT_FOUND_ID, err := utils.UUID()
	assert.Nil(t, err)

	_, err = sut.Answer(fixtures.DUMMY_CONTEXT, c_case.AnswerInput{
		ChallengeID: NOT_FOUND_ID,
		OptionID:    NOT_FOUND_ID,
	})

	assert.ErrorIs(t, err, c_case.ErrOptionNotFound)
}

func Test_List__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()
	challenge := fixtures.CreateChallenge(t)

	sut := fixtures.ChallengeCase.NewIntegration()

	results, err := sut.List()
	require.Nil(t, err)
	require.Len(t, results, 1)

	challengeFound := results[0]
	require.Equal(t, challengeFound.ID, challenge.ID)
	require.Equal(t, challengeFound.Title, challenge.Title)
	require.Equal(t, challengeFound.Description, challenge.Media.Description)
	require.Equal(t, challengeFound.Media, challenge.Media.Path)
}
