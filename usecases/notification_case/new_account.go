package notification_case

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/gateways/notification_gateway"
	"time"
)

type NewAccountInput struct {
	Login    string
	Email    string
	Token    string
	Name     string
	Schedule *time.Time
}

func (n *NotificationCase) NewAccount(
	ctx *context.CoreCtx,
	input NewAccountInput,
) error {
	return n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{input.Email},
		Subject:    "Seja bem-vindo ao Conformity Pro",
		TemplateID: "new-account",
		Variables: map[string]string{
			"name":                input.Name,
			"login":               input.Login,
			"link":                config.AppURL,
			"link_reset_password": config.AppURL + "/reset-password/" + input.Token,
		},
		Schedule: input.Schedule,
	})
}
