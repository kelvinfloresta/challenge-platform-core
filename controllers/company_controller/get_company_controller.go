package company_controller

import (
	"conformity-core/usecases/company_case"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	id := c.Params("company_id")

	data, err := company_case.Singleton.GetById(id)

	if err != nil {
		return err
	}

	if data == nil {
		return c.SendStatus(404)
	}

	return c.JSON(data)
}
