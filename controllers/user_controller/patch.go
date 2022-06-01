package user_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/user_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Patch(c *fiber.Ctx) error {
	data := user_case.PatchInput{}
	if err := utils.ValidateBody(&data, c); err != nil {
		return c.JSON(err)
	}

	data.ID = c.Params("id")
	data.DepartmentID = c.Params("department_id")

	patched, err := user_case.Singleton.Patch(context.New(c.Context()), data)

	if err != nil {
		return err
	}

	if patched {
		return c.SendStatus(200)
	}

	return c.SendStatus(404)
}
