package campaign_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func GetDepartmentAVGResult(c *fiber.Ctx) error {
	data := &campaign_case.GetDepartmentAVGResultInput{}

	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	result, err := campaign_case.Singleton.GetDepartmentAVGResult(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	return c.JSON(result)
}
