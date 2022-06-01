package user_case

import (
	"conformity-core/context"
	"conformity-core/usecases/notification_case"
	"fmt"

	"github.com/getsentry/sentry-go"
)

type NotifyUserWithoutEmailInput struct {
	UserID    string
	CompanyID string
	Name      string
	Phone     string
}

func (u *UserCase) NotifyUserWithoutEmail(input *NotifyUserWithoutEmailInput) {
	hub := sentry.CurrentHub().Clone()
	scope := hub.Scope()
	scope.SetTags(map[string]string{
		"companyId": input.CompanyID,
		"phone":     input.Phone,
	})
	scope.SetUser(sentry.User{
		ID:       input.UserID,
		Username: input.Name,
	})

	company, err := u.companyCase.GetById(input.CompanyID)

	if err != nil {
		hub.CaptureException(fmt.Errorf("GetById: %v", err))
		return
	}

	if company == nil {
		hub.CaptureMessage("Company not found with id: " + input.CompanyID)
		return
	}

	err = u.notificationCase.NewAccountWithoutEmail(context.Internal, notification_case.NewAccountWithoutEmailInput{
		UserID:      input.UserID,
		Name:        input.Name,
		Phone:       input.Phone,
		CompanyID:   input.CompanyID,
		CompanyName: company.Name,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("NewAccountWithoutEmail: %v", err))
		return
	}
}
