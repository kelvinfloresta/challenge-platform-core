package challenge_gateway

import (
	"time"
)

type IChallengeGateway interface {
	Create(input CreateChallengeGatewayInput) (string, error)
	GetById(id string) (*GetByIdChallengeGatewayOutput, error)
	List() ([]ListOutput, error)
	GetOption(input GetOptionInput) (*Option, error)
	GetQuestions(challengeIds []string) (*[]GetQuestionsOutput, error)
}

type ListOutput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Segment     string `json:"segment"`
	Media       string `json:"media"`
}

type CreateChallengeGatewayInput struct {
	Title     string
	Segment   string
	Media     CreateMedia
	Questions []CreateQuestion
}

type GetByIdChallengeGatewayOutput struct {
	ID        string
	Title     string
	Segment   string
	Media     Media
	Questions []Question

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
