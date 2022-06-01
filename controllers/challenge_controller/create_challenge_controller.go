package challenge_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/challenge_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	data := &challenge_case.CreateChallengeCaseInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	id, err := challenge_case.Singleton.Create(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	return c.Status(201).SendString(id)
}
