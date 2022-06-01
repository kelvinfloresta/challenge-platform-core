package campaign_case

import (
	g "conformity-core/gateways/campaign_gateway"
)

func (u CampaignCase) GetById(id string) (*g.GetByIdCampaignGatewayOutput, error) {
	return u.gateway.GetById(id)
}
