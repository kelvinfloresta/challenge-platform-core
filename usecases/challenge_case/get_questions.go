package challenge_case

import g "conformity-core/gateways/challenge_gateway"

func (c ChallengeCase) GetQuestions(challengesIds []string) (*[]g.GetQuestionsOutput, error) {
	return c.gateway.GetQuestions(challengesIds)
}
