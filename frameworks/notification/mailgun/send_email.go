package mailgun

import (
	"conformity-core/config"
	"conformity-core/context"
	"conformity-core/gateways/notification_gateway"
	goContext "context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
)

func SendEmail(ctx *context.CoreCtx, input notification_gateway.SendEmailInput) error {
	mg := mailgun.NewMailgun(config.MailDomain, config.MailKey)

	message := mg.NewMessage(config.MailSender, input.Subject, "", input.To...)

	message.SetTemplate(input.TemplateID)
	if input.Schedule != nil {
		message.SetDeliveryTime(*input.Schedule)
	}

	if input.Attachment != nil {
		message.AddBufferAttachment(input.Attachment.Name, input.Attachment.File)
	}

	for k, v := range input.Variables {
		if err := message.AddTemplateVariable(k, v); err != nil {
			return err
		}
	}

	withTimeout, cancel := goContext.WithTimeout(ctx.Context, time.Second*30)
	defer cancel()

	resp, id, err := mg.Send(withTimeout, message)

	if err != nil {
		return err
	}

	emails := fmt.Sprintf("%s...", input.To[0])

	ctx.Logger.WithFields(logrus.Fields{
		"id":     id,
		"resp":   resp,
		"emails": emails,
	}).Info("mail-sended")

	return nil
}
