package user_case

import (
	coreContext "conformity-core/context"
	"conformity-core/enums/user_company_enum"
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/notification_case"
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

type NotifyNewAccountInput struct {
	DepartmentID string
	UserID       string
	Name         string
	Email        string
	Login        string
	Role         user_company_enum.Role
	Schedule     *time.Time
}

func (u *UserCase) NotifyNewAccount(input *NotifyNewAccountInput) {
	hub := sentry.CurrentHub().Clone()
	scope := hub.Scope()
	scope.SetUser(sentry.User{
		Email: input.Email,
	})

	company, err := u.companyCase.GetOneByFilter(company_case.GetOneByFilterInput{
		DepartmentID: input.DepartmentID,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("GetOneByFilter: %v", err))
		return
	}

	if company == nil {
		hub.CaptureMessage("Company not found with DepartmentID: " + input.DepartmentID)
		return
	}

	if company.RequirePassword {
		u.notifyNewAccount(input, hub)
		return
	}

	u.notifyNewAccountWithoutPassword(input, hub)
}

func (u *UserCase) notifyNewAccount(input *NotifyNewAccountInput, hub *sentry.Hub) {
	token, err := generateResetPasswordToken(input.UserID, input.Schedule)

	if err != nil {
		hub.CaptureException(fmt.Errorf("generateResetPasswordToken: %v", err))
		return
	}

	ctx := coreContext.New(context.Background())
	err = u.notificationCase.NewAccount(ctx, notification_case.NewAccountInput{
		Login:    input.Login,
		Email:    input.Email,
		Token:    token,
		Name:     input.Name,
		Schedule: input.Schedule,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("NewAccount: %v", err))
		return
	}

	if input.Role != user_company_enum.RoleCompanyManager {
		return
	}

	err = u.notificationCase.NewManager(ctx, notification_case.NewManagerInput{
		Email: input.Email,
		Name:  input.Name,
		Token: token,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("manager: %v", err))
	}

}

func (u *UserCase) notifyNewAccountWithoutPassword(input *NotifyNewAccountInput, hub *sentry.Hub) {
	err := u.notificationCase.NewAccountWithoutPassword(notification_case.NewAccountWithoutPasswordInput{
		Email:    input.Email,
		Login:    input.Login,
		Name:     input.Name,
		Schedule: input.Schedule,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("NewAccount: %v", err))
	}

	if input.Role != user_company_enum.RoleCompanyManager {
		return
	}

	token, err := generateResetPasswordToken(input.UserID, input.Schedule)

	if err != nil {
		hub.CaptureException(fmt.Errorf("generateResetPasswordToken: %v", err))
		return
	}

	ctx := coreContext.New(context.Background())
	err = u.notificationCase.NewManager(ctx, notification_case.NewManagerInput{
		Email: input.Email,
		Name:  input.Name,
		Token: token,
	})

	if err != nil {
		hub.CaptureException(fmt.Errorf("manager: %v", err))
	}

}
