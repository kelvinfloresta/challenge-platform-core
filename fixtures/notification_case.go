package fixtures

import (
	coreContext "conformity-core/context"
	"conformity-core/gateways/notification_gateway"
	"conformity-core/usecases/notification_case"
)

type notificationCaseFixture struct {
}

var NotificationCase notificationCaseFixture = notificationCaseFixture{}

type SendEmailMock struct {
	Messages        []notification_gateway.SendEmailInput
	SendEmailOutput error
}

func (s *SendEmailMock) SendEmail(ctx *coreContext.CoreCtx, input notification_gateway.SendEmailInput) error {
	s.Messages = append(s.Messages, input)
	return s.SendEmailOutput
}

func (u *notificationCaseFixture) NewUnit() (*notification_case.NotificationCase, *SendEmailMock) {
	sendEmail := SendEmailMock{}
	return notification_case.New(sendEmail.SendEmail), &sendEmail
}
