package notification_case

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/gateways/notification_gateway"
)

type NewManagerInput struct {
	Email string
	Name  string
	Token string
}

func (n *NotificationCase) NewManager(ctx *context.CoreCtx, input NewManagerInput) error {
	return n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{input.Email},
		Subject:    "Conformity Pro - Gestor",
		TemplateID: "new-manager",
		Variables: map[string]string{
			"name":                input.Name,
			"link":                config.AppURL,
			"link_reset_password": config.AppURL + "/reset-password/" + input.Token,
		},
		Schedule: nil,
	})
}
