package import_controller

import (
	"conformity-core/context"
	"conformity-core/gateways/import_gateway"
	"conformity-core/usecases/import_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func ImportUser(c *fiber.Ctx) error {
	data := &import_gateway.ImportUsersInput{}

	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	err := import_case.Singleton.ImportUsers(context.New(c.Context()), data)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}
