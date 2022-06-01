package routes

import (
	"conformity-core/controllers/company_controller"
	"conformity-core/controllers/department_controller"

	"github.com/gofiber/fiber/v2"
)

func Departments(app *fiber.App) {
	departmentRoutes := app.Group("/departments")
	departmentRoutes.Post("/", department_controller.Create)
	departmentRoutes.Get("/", department_controller.List)
	departmentRoutes.Patch("/:id", department_controller.Patch)
	departmentRoutes.Delete("/:id", department_controller.Delete)
	departmentRoutes.Delete("/:department_id/user/:user_id", company_controller.RemoveUser)
}
