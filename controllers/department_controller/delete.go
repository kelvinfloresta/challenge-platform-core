package department_controller

import (
	"conformity-core/context"
	"conformity-core/usecases/department_case"

	"github.com/gofiber/fiber/v2"
)

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	deleted, err := department_case.Singleton.Delete(context.New(c.Context()), id)

	if err != nil {
		return err
	}

	if deleted {
		return c.SendStatus(200)
	}

	return c.SendStatus(404)
}
