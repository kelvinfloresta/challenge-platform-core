package notification_case

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/gateways/notification_gateway"
)

type ResetPasswordInput struct {
	Email string
	Token string
	Name  string
}

func (n *NotificationCase) ResetPassword(
	ctx *context.CoreCtx,
	input ResetPasswordInput,
) error {
	return n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{input.Email},
		Subject:    "Redefinição de senha",
		TemplateID: "reset-password",
		Variables: map[string]string{
			"name": input.Name,
			"link": config.AppURL + "/reset-password/" + input.Token,
		},
		Schedule: nil,
	})
}
