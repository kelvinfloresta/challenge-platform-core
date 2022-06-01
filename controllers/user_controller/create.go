package user_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/user_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	data := &user_case.CreateInput{}

	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	id, err := user_case.Singleton.Create(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	return c.Status(201).SendString(id)
}
