package company_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/company_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	data := &company_case.CreateCompanyCaseInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	id, err := company_case.Singleton.Create(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	return c.SendString(id)
}
