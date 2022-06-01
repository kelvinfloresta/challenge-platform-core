package company_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/company_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func ChangeUserRole(c *fiber.Ctx) error {
	data := &company_case.ChangeUserRoleCaseInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	success, err := company_case.Singleton.ChangeUserRole(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	if success {
		return c.SendStatus(204)
	}

	return c.SendStatus(404)
}
