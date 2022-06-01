package app

import (
	"conformity-core/config"
	"conformity-core/controllers/auth_controller"
	"conformity-core/controllers/user_controller"
	"conformity-core/routes"

	sentryfiber "github.com/aldy505/sentry-fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
)

func CreateApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(
		sentryfiber.New(sentryfiber.Options{}),
	).Use(
		cors.New(),
	).Use(
		requestid.New(),
	).Use(
		performanceMiddleware,
	).Use(
		logMiddleware,
	).Post(
		"/login", auth_controller.Login,
	).Post(
		"/login/:login", auth_controller.LoginWithoutPassword,
	).Post(
		"/reset-password", user_controller.ResetPassword,
	).Post(
		"/change-password",
		jwtware.New(jwtware.Config{
			SigningKey: []byte(config.ResetPasswordSecret),
		}),
		auth_controller.Session,
		user_controller.ChangePassword,
	).Use(
		"/saml/:workspace", SingleSignOnSAML(app),
	).Use(
		jwtware.New(jwtware.Config{
			SigningKey: []byte(config.AuthSecret),
		}),
	).Use(
		auth_controller.Session,
	).Use("/impersonate", auth_controller.Impersonate)

	routes.Users(app)
	routes.Companies(app)
	routes.Challenges(app)
	routes.Campaigns(app)
	routes.Departments(app)
	routes.Import(app)

	return app
}
