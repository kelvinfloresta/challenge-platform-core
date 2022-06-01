package notification_case

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/frameworks/logger"
	"conformity-core/gateways/notification_gateway"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func (n *NotificationCase) DeadlineChallenge(
	ctx *context.CoreCtx,
	email string,
	name string,
	schedule *time.Time,
) {
	const notificationName = "notify-deadline-challenge"
	if email == "" {
		return
	}

	log := logger.Internal.WithFields(logrus.Fields{
		"notification": notificationName,
		"email":        email,
	})

	err := n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{email},
		Subject:    "Você possui um desafio pendente",
		TemplateID: "deadline-challenge",
		Variables: map[string]string{
			"name": name,
			"link": config.AppURL,
		},
		Schedule: schedule,
	})

	if err != nil {
		log.Error("send-email")
		sentry.CaptureException(err)
	}
}
