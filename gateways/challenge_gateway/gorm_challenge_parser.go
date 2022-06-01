package challenge_gateway

import (
	"conformity-core/frameworks/database/gorm/models"
	"conformity-core/utils"
	"encoding/json"
)

func parseMedia(media CreateMedia) ([]byte, error) {
	var parsed models.ChallengeMediaJSON

	ccId, err := utils.UUID()
	if err != nil {
		return nil, err
	}

	parsed = models.ChallengeMediaJSON{
		ID:          ccId,
		Title:       media.Title,
		Path:        media.Path,
		Description: media.Description,
	}

	return json.Marshal(parsed)
}

func parseQuestions(questions []CreateQuestion) ([]byte, error) {
	var parsed []models.ChallengeQuestionJSON

	for _, question := range questions {
		qId, err := utils.UUID()
		if err != nil {
			return nil, err
		}

		options, err := parseOptions(question.Options)
		if err != nil {
			return nil, err
		}

		parsed = append(parsed, models.ChallengeQuestionJSON{
			ID:      qId,
			Title:   question.Title,
			Options: options,
		})
	}

	return json.Marshal(parsed)
}

func parseOptions(options []CreateOption) ([]models.ChallengeOptionJSON, error) {
	var parsed []models.ChallengeOptionJSON

	for _, option := range options {
		qId, err := utils.UUID()
		if err != nil {
			return nil, err
		}

		parsed = append(parsed, models.ChallengeOptionJSON{
			ID:      qId,
			Title:   option.Title,
			Correct: option.Correct,
		})
	}

	return parsed, nil
}
