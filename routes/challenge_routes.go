package routes

import (
	"conformity-core/controllers/challenge_controller"

	"github.com/gofiber/fiber/v2"
)

func Challenges(app *fiber.App) {
	challengeRoutes := app.Group("/challenges")
	challengeRoutes.Post("/", challenge_controller.Create)
	challengeRoutes.Get("/", challenge_controller.List)
	challengeRoutes.Get("/:challenge_id", challenge_controller.Get)
}
