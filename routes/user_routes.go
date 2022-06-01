package routes

import (
	"conformity-core/controllers/user_controller"

	"github.com/gofiber/fiber/v2"
)

func Users(app *fiber.App) {
	userRoutes := app.Group("/users")
	userRoutes.Post("/", user_controller.Create)
	userRoutes.Patch("/:id/:department_id", user_controller.Patch)
}
