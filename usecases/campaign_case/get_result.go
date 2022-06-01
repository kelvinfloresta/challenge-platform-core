package campaign_case

import (
	"conformity-core/context"
	g "conformity-core/gateways/campaign_gateway"
)

func (u CampaignCase) GetResult(ctx *context.CoreCtx, input g.GetResultInput) (*g.Result, error) {
	return u.gateway.GetResult(ctx, input)
}
