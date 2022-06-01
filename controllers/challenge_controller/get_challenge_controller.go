package challenge_controller

import (
	"conformity-core/usecases/challenge_case"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	id := c.Params("challenge_id")

	data, err := challenge_case.Singleton.GetById(id)

	if err != nil {
		return err
	}

	if data == nil {
		return c.SendStatus(404)
	}

	return c.JSON(data)
}
