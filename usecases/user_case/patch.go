package user_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/gateways/user_gateway"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type PatchInput struct {
	ID           string
	OID          string
	DepartmentID string
	Name         string
	Document     string
	Email        string `json:"email" validate:"omitempty,email"`
	Phone        string `json:"phone" validate:"omitempty,e164"`
	JobPosition  string
	Role         user_company_enum.Role `validate:"oneof=user companyManager"`
}

func (u UserCase) Patch(ctx *context.CoreCtx, input PatchInput) (bool, error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager &&
		ctx.Session.Role != user_company_enum.RoleBackoffice {
		return false, core_errors.ErrWithoutPrivileges
	}

	shouldNotifyNewPhone, err := u.shouldNotifyBackofficeAboutNewPhone(input)

	if err != nil {
		return false, err
	}

	updated, err := u.gateway.Patch(user_gateway.PatchInput{
		ID:           input.ID,
		OID:          input.OID,
		DepartmentID: input.DepartmentID,
		Email:        input.Email,
		Document:     input.Document,
		Login:        getLogin(input),
		Name:         input.Name,
		Phone:        input.Phone,
		JobPosition:  input.JobPosition,
		Role:         input.Role,
	})

	if err != nil {
		return false, err
	}

	if shouldNotifyNewPhone {
		go u.NotifyUserWithoutEmail(&NotifyUserWithoutEmailInput{
			UserID:    input.ID,
			CompanyID: ctx.Session.CompanyID,
			Name:      input.Name,
			Phone:     input.Phone,
		})
	}

	ctx.Logger.WithFields(logrus.Fields{
		"id":           input.ID,
		"oid":          input.OID,
		"departmentId": input.DepartmentID,
		"email":        input.Email,
		"document":     input.Document,
		"login":        getLogin(input),
		"name":         input.Name,
		"phone":        input.Phone,
		"jobPosition":  input.JobPosition,
		"role":         input.Role,
		"updated":      updated,
	}).Info("update-user-department")

	return updated, nil
}

func getLogin(input PatchInput) string {
	if input.Email != "" && input.Document != "" {
		return ""
	}

	if input.Email != "" {
		return input.Email
	}

	return input.Document
}

func (u UserCase) shouldNotifyBackofficeAboutNewPhone(input PatchInput) (bool, error) {
	userReceiveNotificationWithEmail := input.Email != ""
	if userReceiveNotificationWithEmail {
		return false, nil
	}

	user, err := u.gateway.GetOneByFilter(user_gateway.GetOneByFilterInput{
		ID:           input.ID,
		DepartmentID: input.DepartmentID,
	})

	if err != nil {
		return false, err
	}

	if user == nil {
		message := fmt.Sprintf(
			"User not found: %s department: %s",
			input.ID,
			input.DepartmentID,
		)
		sentry.CaptureMessage(message)
		return false, nil
	}

	userWithNewPhone := user.Phone != input.Phone
	return userWithNewPhone, nil
}
