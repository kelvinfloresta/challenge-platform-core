package auth_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/auth_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Impersonate(c *fiber.Ctx) error {
	data := auth_case.ImpersonateInput{}

	if err := utils.ValidateBody(&data, c); err != nil {
		return c.JSON(err)
	}

	token, err := auth_case.Singleton.Impersonate(context.New(c.Context()), data)
	if err != nil {
		return err
	}

	if token == "" {
		return c.SendStatus(404)
	}

	return c.SendString(token)
}
