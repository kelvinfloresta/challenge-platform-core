package campaign_case

import (
	"conformity-core/context"
	"conformity-core/domain/challenge"
	"conformity-core/enums/campaign_enum"
	g "conformity-core/gateways/campaign_gateway"
	"math/rand"
	"time"
)

type ListQuestionsInput struct {
	ChallengeID string `validate:"required" json:"challenge_id"`
	CampaignID  string `validate:"required" json:"campaign_id"`
}

type ListOptions struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ListQuestionsOutput struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Options        []ListOptions `json:"options"`
	RemainingTries uint8         `json:"remainingTries"`
}

func (u CampaignCase) ListQuestions(ctx *context.CoreCtx, input ListQuestionsInput) ([]ListQuestionsOutput, error) {
	var triesLessThan uint8 = 3
	correct := false

	data, err := u.gateway.ListQuestions(ctx, g.ListQuestionsInput{
		ChallengeID:   input.ChallengeID,
		CampaignID:    input.CampaignID,
		UserID:        ctx.Session.UserID,
		DepartmentID:  ctx.Session.DepartmentID,
		TriesLessThan: &triesLessThan,
		Correct:       &correct,
		Status:        campaign_enum.Active,
	})

	if err != nil {
		return nil, err
	}

	result := parseListQuestions(data)

	return result, nil
}

func parseListQuestions(data []g.ListQuestionsOutput) []ListQuestionsOutput {
	result := []ListQuestionsOutput{}

	for _, lqo := range data {
		options := make([]ListOptions, 0, len(lqo.Options))
		for _, o := range lqo.Options {
			options = append(options, ListOptions{
				ID:    o.ID,
				Title: o.Title,
			})
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })

		result = append(result, ListQuestionsOutput{
			ID:             lqo.ID,
			Title:          lqo.Title,
			Options:        options,
			RemainingTries: challenge.MAX_TRIES - lqo.Tries,
		})
	}

	return result
}
