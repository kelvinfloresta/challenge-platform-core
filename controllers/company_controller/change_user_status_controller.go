package company_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func ChangeUserStatus(c *fiber.Ctx) error {
	data := &campaign_case.ChangeUserStatusCaseInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	success, err := campaign_case.Singleton.ChangeUserStatus(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	if success {
		return c.SendStatus(204)
	}

	return c.SendStatus(404)
}
