package campaign_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"

	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	data := campaign_case.ListUsersInput{}
	if err := c.QueryParser(&data); err != nil {
		return err
	}

	users, err := campaign_case.Singleton.ListUsers(
		context.New(c.Context()),
		data,
	)

	if err != nil {
		return err
	}

	return c.JSON(users)
}
