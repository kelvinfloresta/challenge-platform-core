package department_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/department_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Patch(c *fiber.Ctx) error {
	data := department_case.PatchInput{}
	if err := utils.ValidateBody(&data, c); err != nil {
		return c.JSON(err)
	}

	data.ID = c.Params("id")

	deleted, err := department_case.Singleton.Patch(context.New(c.Context()), data)

	if err != nil {
		return err
	}

	if deleted {
		return c.SendStatus(200)
	}

	return c.SendStatus(404)
}
