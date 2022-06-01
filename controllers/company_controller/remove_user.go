package company_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"

	"github.com/gofiber/fiber/v2"
)

func RemoveUser(c *fiber.Ctx) error {
	data := &campaign_case.RemoveUserInput{
		UserID:       c.Params("user_id"),
		DepartmentID: c.Params("department_id"),
	}

	removed, err := campaign_case.Singleton.RemoveUser(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	if removed {
		return c.SendStatus(201)
	}

	return c.SendStatus(404)

}
