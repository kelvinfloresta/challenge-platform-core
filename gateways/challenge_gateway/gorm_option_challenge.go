package challenge_gateway

import "conformity-core/frameworks/database/gorm/models"

func (g *GormChallengeGatewayFacade) GetOption(input GetOptionInput) (*Option, error) {
	data := &models.ChallengeOptionJSON{}

	filter := map[string]interface{}{
		"challengeId": input.ChallengeID,
		"questionId":  input.QuestionID,
		"optionId":    input.OptionID,
	}

	result := g.DB.Conn.Raw(`
SELECT
	target."Options"->>'ID' as ID,
	target."Options"->>'Title' as Title,
	target."Options"->'Correct' as Correct
FROM
	"challenges",
	jsonb_array_elements("challenges".questions) dependency("Questions"),
	jsonb_array_elements("Questions"->'Options') target("Options")
WHERE
	"challenges".id = @challengeId
AND
	dependency."Questions"->>'ID' = @questionId
AND 
	target."Options"->>'ID' = @optionId`, filter).Scan(data)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := &Option{
		ID:      data.ID,
		Title:   data.Title,
		Correct: data.Correct,
	}

	return adapted, nil
}
