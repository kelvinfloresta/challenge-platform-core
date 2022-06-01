package campaign_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func ListResults(c *fiber.Ctx) error {
	data := campaign_case.ListResultsInput{}

	if err := utils.ValidateBody(&data, c); err != nil {
		return c.JSON(err)
	}

	results, err := campaign_case.Singleton.ListResults(context.New(c.Context()), data)

	if err != nil {
		return err
	}

	return c.JSON(results)
}
