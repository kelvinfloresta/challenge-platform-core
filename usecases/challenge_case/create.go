package challenge_case

import (
	"conformity-core/context"
	g "conformity-core/gateways/challenge_gateway"
)

type CreateChallengeCaseInput struct {
	Title     string             `json:"title" validate:"required"`
	Segment   string             `json:"segment" validate:"required"`
	Media     g.CreateMedia      `json:"media" validate:"required"`
	Questions []g.CreateQuestion `json:"questions" validate:"required"`
}

func (c ChallengeCase) Create(ctx *context.CoreCtx, input CreateChallengeCaseInput) (id string, err error) {
	err = validateChallenge(input)
	if err != nil {
		return
	}

	id, err = c.gateway.Create(g.CreateChallengeGatewayInput{
		Title:     input.Title,
		Segment:   input.Segment,
		Media:     input.Media,
		Questions: input.Questions,
	})

	if err != nil {
		return
	}

	return
}

func validateOptions(options []g.CreateOption) (err error) {
	var corrects int
	for _, opt := range options {
		if opt.Correct {
			corrects++
		}

		if corrects > 1 {
			return ErrMultipleCorrectOptions
		}
	}

	if corrects == 0 {
		return ErrWithoutCorrectOption
	}

	return
}

func validateChallenge(input CreateChallengeCaseInput) (err error) {
	if len(input.Questions) == 0 {
		return ErrEmptyChallengeQuestions
	}

	for _, cq := range input.Questions {
		if err = validateOptions(cq.Options); err != nil {
			return
		}
	}

	return
}
