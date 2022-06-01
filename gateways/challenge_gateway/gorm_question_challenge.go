package challenge_gateway

func (g *GormChallengeGatewayFacade) GetQuestions(challengeIds []string) (*[]GetQuestionsOutput, error) {
	output := &[]GetQuestionsOutput{}
	result := g.DB.Conn.Raw(`
		SELECT
			"challenges".id as challenge_id,
			arr."question"->>'ID' as question_id
		FROM
			"challenges",
			jsonb_array_elements("challenges".questions) arr("question")
		WHERE "challenges".id IN (?)`, challengeIds).Scan(output)

	if result.Error != nil {
		return nil, result.Error
	}

	return output, nil
}
