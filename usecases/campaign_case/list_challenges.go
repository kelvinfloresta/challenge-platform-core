package campaign_case

import (
	"conformity-core/context"
	"conformity-core/domain/challenge"
	"conformity-core/enums/campaign_enum"
	g "conformity-core/gateways/campaign_gateway"
	"time"
)

type ListChallengesInput struct {
	IsPending *bool
	UserID    string
}

func (c CampaignCase) ListChallenges(
	ctx *context.CoreCtx,
	input ListChallengesInput,
) ([]g.GetChallengesCampaignOutput, error) {
	now := time.Now()
	maxTries := challenge.MAX_TRIES

	filter := g.ListChallengesInput{
		Status: campaign_enum.Active,
	}

	if !ctx.Session.Role.IsBackoffice() {
		filter.CompanyID = ctx.Session.CompanyID
	}

	if ctx.Session.Role.IsCommonUser() {
		filter.UserID = ctx.Session.UserID
	} else {
		filter.UserID = input.UserID
	}

	if input.IsPending == nil {
		return c.gateway.ListChallenges(ctx, filter)
	} else if *input.IsPending {
		notCorrect := false
		filter.Correct = &notCorrect
		filter.TriesLT = &maxTries
		filter.EndDateGTE = now
		filter.StartDateLTE = now
	} else {
		correct := true
		filter.OR = &g.ListChallengesORInput{
			Correct:    &correct,
			Tries:      &maxTries,
			EndDateLTE: now,
		}
	}

	return c.gateway.ListChallenges(ctx, filter)
}
