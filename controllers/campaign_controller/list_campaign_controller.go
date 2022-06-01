package campaign_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"

	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	campaigns, err := campaign_case.Singleton.List(context.New(c.Context()))

	if err != nil {
		return err
	}

	return c.JSON(campaigns)
}
