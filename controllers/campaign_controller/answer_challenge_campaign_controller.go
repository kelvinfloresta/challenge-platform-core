package campaign_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func AnswerChallenge(c *fiber.Ctx) error {
	data := &campaign_case.AnswerChallengeInput{}

	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	correct, err := campaign_case.Singleton.AnswerChallenge(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	return c.Status(201).JSON(correct)
}
