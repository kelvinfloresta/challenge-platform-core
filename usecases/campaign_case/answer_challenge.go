package campaign_case

import (
	"conformity-core/context"
	"conformity-core/domain/challenge"
	g "conformity-core/gateways/campaign_gateway"
	"conformity-core/usecases/challenge_case"

	"github.com/sirupsen/logrus"
)

type AnswerChallengeInput struct {
	CampaignID  string `json:"campaign_id" validate:"required"`
	ChallengeID string `json:"challenge_id" validate:"required"`
	QuestionID  string `json:"question_id" validate:"required"`
	OptionID    string `json:"option_id" validate:"required"`
}

type AnswerChallengeOutput struct {
	Correct        bool
	RemainingTries uint8
}

func (c CampaignCase) AnswerChallenge(ctx *context.CoreCtx, input AnswerChallengeInput) (output *AnswerChallengeOutput, err error) {
	userSessionData := ctx.Session
	userId := userSessionData.UserID
	departmentId := userSessionData.DepartmentID

	found, err := c.gateway.GetResult(ctx, g.GetResultInput{
		UserID:       userId,
		DepartmentID: departmentId,
		ChallengeID:  input.ChallengeID,
		CampaignID:   input.CampaignID,
		QuestionID:   input.QuestionID,
	})

	if err != nil {
		return
	}

	if found == nil {
		return nil, ErrUserIsNotInCampaign
	}

	if found.Tries == 3 {
		return nil, ErrNoTriesRemaining
	}

	if found.Correct {
		return nil, ErrQuestionIsAlreadyCorrect
	}

	correct, err := c.challengeCase.Answer(ctx, challenge_case.AnswerInput{
		ChallengeID: input.ChallengeID,
		QuestionID:  input.QuestionID,
		OptionID:    input.OptionID,
	})

	if err != nil {
		return
	}

	newTries := found.Tries + uint8(1)
	updatedResult := g.UpdateResultInput{
		NewTries:     &newTries,
		ActualTries:  &found.Tries,
		NewCorrect:   &correct,
		UserID:       userId,
		DepartmentID: departmentId,
		CampaignID:   input.CampaignID,
		ChallengeID:  input.ChallengeID,
		QuestionID:   input.QuestionID,
	}

	affected, err := c.gateway.UpdateResult(updatedResult, nil)

	if err != nil {
		return
	}

	if !affected {
		ctx.Logger.WithFields(logrus.Fields{
			"CampaignID":  input.CampaignID,
			"OptionID":    input.OptionID,
			"QuestionID":  input.QuestionID,
			"ChallengeID": input.ChallengeID,
		}).Error(ErrFailUpdateCampaignResult.Error())

		return nil, ErrFailUpdateCampaignResult
	}

	newAmountOfTries := found.Tries + 1
	remainingTries := challenge.MAX_TRIES - newAmountOfTries

	output = &AnswerChallengeOutput{
		Correct:        correct,
		RemainingTries: remainingTries,
	}

	return
}
