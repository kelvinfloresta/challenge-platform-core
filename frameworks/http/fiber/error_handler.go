package app

import (
	"conformity-core/context"
	core_errors "conformity-core/errors"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func errorHandler(c *fiber.Ctx, err error) error {
	switch v := err.(type) {
	case core_errors.Unauthorized:
		return c.Status(401).SendString(v.Error())

	case core_errors.Forbidden:
		return c.Status(403).SendString(v.Error())

	case core_errors.NotFound:
		return c.Status(404).SendString(v.Error())

	case core_errors.Conflict:
		return c.Status(409).SendString(v.Error())

	case core_errors.BadRequest:
		return c.Status(400).SendString(v.Error())

	case *fiber.Error:
		return c.Status(v.Code).JSON(err)

	default:
		scope := getOrCreateHub(c).Scope()
		session := context.New(c.Context())
		scope.SetUser(sentry.User{
			ID: session.Session.UserID,
		})
		sentry.CaptureException(err)
		return c.Status(500).SendString("Internal Server Error")
	}
}
