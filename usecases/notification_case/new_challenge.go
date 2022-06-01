package notification_case

import (
	"conformity-core/config"
	coreContext "conformity-core/context"
	"conformity-core/frameworks/logger"
	"conformity-core/gateways/notification_gateway"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

const newChallenge = "notify-start-challenge"

type NewChallengeUser struct {
	Email string
	Name  string
}
type NewChallengeInput struct {
	Users    []NewChallengeUser
	Schedule *time.Time
}

func (n *NotificationCase) NewChallenge(
	ctx *coreContext.CoreCtx,
	input NewChallengeInput,
) {
	for _, user := range input.Users {
		if user.Email == "" {
			continue
		}

		log := logger.Internal.WithFields(logrus.Fields{
			"notification": newChallenge,
			"email":        user.Email,
		})

		err := n.sendEmail(ctx, notification_gateway.SendEmailInput{
			To:         []string{user.Email},
			Subject:    "Novo desafio enviado para você",
			TemplateID: "new-challenge",
			Variables: map[string]string{
				"name": user.Name,
				"link": config.AppURL,
			},
			Schedule: input.Schedule,
		})

		if err != nil {
			log.Error("send-email")
			sentry.CaptureException(err)
			continue
		}
	}
}
