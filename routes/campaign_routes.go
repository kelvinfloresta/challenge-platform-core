package routes

import (
	"conformity-core/controllers/campaign_controller"

	"github.com/gofiber/fiber/v2"
)

func Campaigns(app *fiber.App) {
	campaignRoutes := app.Group("/campaigns")
	campaignRoutes.Post("/", campaign_controller.Create)
	campaignRoutes.Get("/", campaign_controller.List)
	campaignRoutes.Post("/result", campaign_controller.ListResults)
	campaignRoutes.Get("/challenges", campaign_controller.ListChallenges)
	campaignRoutes.Get("/users", campaign_controller.ListUsers)
	campaignRoutes.Post("/challenges/answer", campaign_controller.AnswerChallenge)
	campaignRoutes.Post("/questions", campaign_controller.ListQuestions)
	campaignRoutes.Post("/department/result", campaign_controller.GetDepartmentAVGResult)
	campaignRoutes.Post("/user/result", campaign_controller.GetUserAVGResult)
	campaignRoutes.Get("/:campaign_id", campaign_controller.Get)
}
