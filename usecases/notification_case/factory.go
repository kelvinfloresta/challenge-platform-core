package notification_case

import (
	"conformity-core/frameworks/notification/mailgun"
)

var Singleton = build()

func build() *NotificationCase {
	return New(mailgun.SendEmail)
}
