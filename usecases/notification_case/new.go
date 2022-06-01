package notification_case

import (
	"conformity-core/gateways/notification_gateway"
)

type NotificationCase struct {
	sendEmail notification_gateway.SendEmail
}

func New(
	sendEmail notification_gateway.SendEmail,
) *NotificationCase {
	return &NotificationCase{sendEmail}
}
