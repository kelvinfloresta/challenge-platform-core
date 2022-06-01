package challenge_gateway

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"conformity-core/utils"
	"encoding/json"
)

type GormChallengeGatewayFacade struct {
	DB *database.Database
}

func (g *GormChallengeGatewayFacade) Create(input CreateChallengeGatewayInput) (string, error) {
	challengeID, err := utils.UUID()
	if err != nil {
		return "", err
	}

	media, err := parseMedia(input.Media)
	if err != nil {
		return "", err
	}

	questions, err := parseQuestions(input.Questions)
	if err != nil {
		return "", err
	}

	err = g.DB.Conn.Create(&models.Challenge{
		ID:        challengeID,
		Title:     input.Title,
		Segment:   input.Segment,
		Media:     media,
		Questions: questions,
	}).Error

	if err != nil {
		return "", err
	}

	return challengeID, nil
}

func (g *GormChallengeGatewayFacade) GetById(id string) (res *GetByIdChallengeGatewayOutput, err error) {
	data := &models.Challenge{ID: id}
	err = g.DB.Conn.Find(data).Error

	if err != nil {
		return
	}

	var media Media
	err = json.Unmarshal(data.Media, &media)
	if err != nil {
		return
	}

	var questions []Question
	err = json.Unmarshal(data.Questions, &questions)
	if err != nil {
		return
	}

	res = &GetByIdChallengeGatewayOutput{
		ID:        data.ID,
		Title:     data.Title,
		Segment:   data.Segment,
		Media:     media,
		Questions: questions,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.Time,
	}

	return
}
