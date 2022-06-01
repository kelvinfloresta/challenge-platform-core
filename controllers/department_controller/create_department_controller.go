package department_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/department_case"
	"conformity-core/utils"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	data := &department_case.CreateDepartmentCaseInput{}
	if err := utils.ValidateBody(data, c); err != nil {
		return c.JSON(err)
	}

	id, err := department_case.Singleton.Create(context.New(c.Context()), *data)

	if err != nil {
		return err
	}

	return c.SendString(id)
}
