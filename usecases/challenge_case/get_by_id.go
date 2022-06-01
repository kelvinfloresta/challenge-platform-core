package challenge_case

import (
	g "conformity-core/gateways/challenge_gateway"
)

func (c ChallengeCase) GetById(id string) (*g.GetByIdChallengeGatewayOutput, error) {
	return c.gateway.GetById(id)
}
