package campaign_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"

	"github.com/gofiber/fiber/v2"
)

func ListChallenges(c *fiber.Ctx) error {
	data := campaign_case.ListChallengesInput{}
	if err := c.QueryParser(&data); err != nil {
		return err
	}

	challenges, err := campaign_case.Singleton.ListChallenges(
		context.New(c.Context()),
		data,
	)

	if err != nil {
		return err
	}

	return c.JSON(challenges)
}
