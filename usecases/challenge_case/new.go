package challenge_case

import g "conformity-core/gateways/challenge_gateway"

type ChallengeCase struct {
	gateway g.IChallengeGateway
}

func New(gateway g.IChallengeGateway) *ChallengeCase {
	return &ChallengeCase{gateway}
}
