package user_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/user_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func ResetPassword(c *fiber.Ctx) error {
	data := &user_case.ResetPasswordInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	err := user_case.Singleton.ResetPassword(context.New(c.Context()), data)
	if err != nil {
		return err
	}

	return c.SendStatus(200)
}
