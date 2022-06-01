package challenge_case

import "conformity-core/gateways/challenge_gateway"

type ListOutput struct {
	ID          string
	Title       string
	Description string
	Media       string
	Segment     string
}

func (c ChallengeCase) List() ([]challenge_gateway.ListOutput, error) {
	return c.gateway.List()
}
