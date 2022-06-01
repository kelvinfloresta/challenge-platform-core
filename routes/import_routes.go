package routes

import (
	"conformity-core/controllers/import_controller"

	"github.com/gofiber/fiber/v2"
)

func Import(app *fiber.App) {
	importRoutes := app.Group("/import")
	importRoutes.Use("/users", import_controller.ImportUser)
}
