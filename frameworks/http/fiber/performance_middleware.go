package app

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func performanceMiddleware(c *fiber.Ctx) error {
	hub := getOrCreateHub(c)
	ctx := context.WithValue(c.Context(), sentry.HubContextKey, hub)
	span := sentry.StartSpan(ctx, "request", func(s *sentry.Span) {
		scope := hub.Scope()
		transaction := fmt.Sprintf("%s %s", c.Method(), c.Path())
		scope.SetTransaction(transaction)

		if id, ok := c.Locals("requestid").(string); ok {
			scope.SetTag("requestid", id)
		}
	})

	defer span.Finish()

	return c.Next()
}
