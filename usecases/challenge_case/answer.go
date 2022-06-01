package challenge_case

import (
	"conformity-core/context"
	g "conformity-core/gateways/challenge_gateway"

	"github.com/sirupsen/logrus"
)

type AnswerInput struct {
	ChallengeID string
	QuestionID  string
	OptionID    string
}

func (c ChallengeCase) Answer(ctx *context.CoreCtx, input AnswerInput) (bool, error) {
	option, err := c.gateway.GetOption(g.GetOptionInput(input))
	if err != nil {
		return false, err
	}

	if option == nil {
		ctx.Logger.WithFields(logrus.Fields{
			"ChallengeID": input.ChallengeID,
			"OptionID":    input.OptionID,
		}).Warn(ErrOptionNotFound.Error())
		return false, ErrOptionNotFound
	}

	ctx.Logger.WithFields(logrus.Fields{
		"ChallengeID": input.ChallengeID,
		"OptionID":    input.OptionID,
		"Correct":     option.Correct,
	}).Info("Answer")

	return option.Correct, nil
}
