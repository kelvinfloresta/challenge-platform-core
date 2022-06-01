package challenge_controller

import (
	"conformity-core/usecases/challenge_case"

	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	data, err := challenge_case.Singleton.List()

	if err != nil {
		return err
	}

	return c.JSON(data)
}
