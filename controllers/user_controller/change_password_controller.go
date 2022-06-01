package user_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/user_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func ChangePassword(c *fiber.Ctx) error {
	data := &user_case.ChangePasswordInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	changed, err := user_case.Singleton.ChangePassword(
		context.New(c.Context()),
		data,
	)

	if err != nil {
		return err
	}

	if !changed {
		return c.SendStatus(404)
	}

	return c.SendStatus(204)
}
