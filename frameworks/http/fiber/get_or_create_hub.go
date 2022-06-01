package app

import (
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func getOrCreateHub(c *fiber.Ctx) *sentry.Hub {
	// This hub was created by lib: github.com/aldy505/sentry-fiber
	if v, ok := c.Locals("sentry").(*sentry.Hub); ok {
		return v
	}

	hub := sentry.CurrentHub().Clone()
	c.Locals("sentry", hub)

	return hub
}
