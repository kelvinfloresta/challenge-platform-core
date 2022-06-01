package auth_controller

import (
	"conformity-core/usecases/auth_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	data := auth_case.LoginInput{}
	if err := utils.ValidateBody(&data, c); err != nil {
		return c.JSON(err)
	}

	token, err := auth_case.Singleton.Login(data)
	if err != nil {
		return err
	}

	if token == "" {
		return c.SendStatus(401)
	}

	return c.JSON(fiber.Map{"token": token})
}
