package routes

import (
	"conformity-core/controllers/company_controller"
	"conformity-core/controllers/department_controller"

	"github.com/gofiber/fiber/v2"
)

func Companies(app *fiber.App) {
	companyRoutes := app.Group("/companies")
	companyRoutes.Post("/", company_controller.Create)
	companyRoutes.Get("/users/:actualPage/:pageSize", department_controller.PaginateUsers)
	companyRoutes.Patch("/users/status", company_controller.ChangeUserStatus)
	companyRoutes.Patch("/users/role", company_controller.ChangeUserRole)
}
