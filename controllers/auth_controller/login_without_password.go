package auth_controller

import (
	"conformity-core/usecases/auth_case"

	"github.com/gofiber/fiber/v2"
)

func LoginWithoutPassword(c *fiber.Ctx) error {
	login := c.Params("login")

	token, err := auth_case.Singleton.LoginWithoutPassword(login)

	if err != nil {
		return err
	}

	if token == "" {
		return c.SendStatus(401)
	}

	return c.JSON(fiber.Map{"token": token})
}
