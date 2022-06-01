package notification_case

import (
	"conformity-core/config"
	coreContext "conformity-core/context"
	"conformity-core/gateways/notification_gateway"
	"context"
	"time"
)

type NewAccountWithoutPasswordInput struct {
	Email    string
	Login    string
	Name     string
	Schedule *time.Time
}

func (n *NotificationCase) NewAccountWithoutPassword(
	input NewAccountWithoutPasswordInput,
) error {
	ctx := coreContext.New(context.Background())
	return n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{input.Email},
		Subject:    "Seja bem-vindo ao Conformity Pro",
		TemplateID: "new-account-without-password",
		Variables: map[string]string{
			"name":  input.Name,
			"login": input.Login,
			"link":  config.AppURL,
		},
		Schedule: input.Schedule,
	})
}
