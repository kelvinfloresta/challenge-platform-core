package challenge_case

import "errors"

var (
	ErrMultipleCorrectOptions  = errors.New("MULTIPLE_CORRECT_OPTIONS")
	ErrWithoutCorrectOption    = errors.New("WITHOUT_CORRECT_OPTION")
	ErrEmptyChallengeMedia     = errors.New("EMPTY_CHALLENGE_MEDIA")
	ErrEmptyChallengeQuestions = errors.New("EMPTY_CHALLENGE_QUESTIONS")
	ErrOptionNotFound          = errors.New("OPTION_NOT_FOUND")
)
