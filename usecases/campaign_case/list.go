package campaign_case

import (
	"conformity-core/context"
	"conformity-core/gateways/campaign_gateway"
)

func (c CampaignCase) List(ctx *context.CoreCtx) ([]*campaign_gateway.ListOutput, error) {
	companyId := ctx.Session.CompanyID

	return c.gateway.List(ctx, campaign_gateway.ListInput{
		CompanyID: companyId,
	})
}
