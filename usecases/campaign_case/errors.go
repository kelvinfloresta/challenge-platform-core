package campaign_case

import (
	"errors"
	"fmt"
)

var ErrCampaignWithoutParticipants = errors.New("CAMPAIGN_WITHOUT_PARTICIPANTS")
var ErrFailUpdateCampaignResult = errors.New("FAIL_UPDATE_CAMPAIGN_RESULT")
var ErrUserIsNotInCampaign = errors.New("USER_IS_NOT_IN_CAMPAIGN")
var ErrNoTriesRemaining = errors.New("NO_TRIES_REMAINING")
var ErrQuestionIsAlreadyCorrect = errors.New("QUESTION_IS_ALREADY_CORRECT")
var ErrCampaignWithoutChallenges = errors.New("CAMPAIGN_WITHOUT_CHALLENGES")
var ErrCampaignWithoutQuestions = errors.New("CAMPAIGN_WITHOUT_QUESTIONS")
var ErrCampaignWithDuplicatedChallengesIds = errors.New("CAMPAIGN_WITH_DUPLICATED_CHALLENGES_IDS")

func ErrPastEndDate(index int) error {
	return errors.New(fmt.Sprint("PAST_END_DATE__", index))
}

func ErrStartDateBeforeToday(index int) error {
	return errors.New(fmt.Sprint("START_DATE_BEFORE_TODAY__", index))
}

func ErrInvalidRange(index int) error {
	return errors.New(fmt.Sprint("INVALID_RANGE__", index))
}
