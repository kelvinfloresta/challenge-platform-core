package notification_case

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/gateways/notification_gateway"
	"fmt"
)

type NewAccountWithoutEmailInput struct {
	UserID      string
	Name        string
	Phone       string
	CompanyID   string
	CompanyName string
}

func (n *NotificationCase) NewAccountWithoutEmail(ctx *context.CoreCtx, input NewAccountWithoutEmailInput) error {
	return n.sendEmail(ctx, notification_gateway.SendEmailInput{
		To:         []string{config.MailSender},
		Subject:    fmt.Sprintf("Novo usuário - %s", input.Phone),
		TemplateID: "new-account-without-email",
		Variables: map[string]string{
			"name":        input.Name,
			"userId":      input.UserID,
			"phone":       input.Phone,
			"companyId":   input.CompanyID,
			"companyName": input.CompanyName,
		},
		Schedule: nil,
	})
}
