package department_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/department_case"

	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	departments, err := department_case.Singleton.List(context.New(c.Context()))

	if err != nil {
		return err
	}

	return c.JSON(departments)
}
