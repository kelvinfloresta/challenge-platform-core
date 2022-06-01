package user_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/gateways/user_gateway"
	"conformity-core/usecases/campaign_case"
	"conformity-core/utils"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type CreateInput struct {
	OID          string                 `json:"-"`
	Name         string                 `json:"name" validate:"required"`
	Password     string                 `json:"password" validate:"omitempty,min=8"`
	Email        string                 `json:"email" validate:"omitempty,email"`
	Login        string                 `json:"login"`
	DepartmentID string                 `json:"departmentId" validate:"required"`
	Document     string                 `json:"document"`
	Phone        string                 `json:"phone" validate:"omitempty,e164"`
	BirthDate    string                 `json:"birthDate"`
	JobPosition  string                 `json:"jobPosition"`
	Role         user_company_enum.Role `json:"role" validate:"required,oneof=user companyManager"`
}

func (c *UserCase) Create(
	ctx *context.CoreCtx,
	input CreateInput,
) (string, error) {
	login, err := validateCreate(ctx, input)
	if err != nil {
		return "", err
	}

	input.Login = login

	if input.Password == "" {
		password, err := utils.UUID()
		if err != nil {
			return "", err
		}

		input.Password = password
	}

	encryptedPassword, err := utils.PasswordGen(input.Password)
	if err != nil {
		return "", err
	}

	input.Password = encryptedPassword

	id, err := c.gateway.Create(user_gateway.CreateInput(input))
	if err != nil {
		return id, err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"name":      input.Name,
		"email":     input.Email,
		"create-id": id,
	}).Info("user-created")

	err = c.campaignCase.AddUserToCampaigns(ctx, campaign_case.AddUserToCampaignsInput{
		UserID:       id,
		DepartmentID: input.DepartmentID,
		CompanyID:    ctx.Session.CompanyID,
	})

	if err != nil {
		hub := sentry.CurrentHub()
		hub.Scope().SetUser(sentry.User{
			ID:       id,
			Email:    input.Email,
			Username: input.Name,
		})
		sentry.CaptureException(fmt.Errorf("cannot add user to campaigns: %v", err))
	}

	if input.Email != "" {
		go c.NotifyNewAccount(&NotifyNewAccountInput{
			UserID:       id,
			DepartmentID: input.DepartmentID,
			Email:        input.Email,
			Name:         input.Name,
			Login:        input.Login,
		})
	}

	if input.Document == login {
		go c.NotifyUserWithoutEmail(&NotifyUserWithoutEmailInput{
			UserID:    id,
			CompanyID: ctx.Session.CompanyID,
			Name:      input.Name,
			Phone:     input.Phone,
		})
	}

	return id, nil
}

func validateCreate(ctx *context.CoreCtx, input CreateInput) (string, error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager &&
		ctx.Session.Role != user_company_enum.RoleBackoffice {
		return "", core_errors.ErrWithoutPrivileges
	}

	if input.Login != "" {
		return input.Login, nil
	}

	if input.Email != "" {
		return input.Email, nil
	}

	if input.Document == "" {
		return "", core_errors.ErrEmailAndDocumentEmpty
	}

	if input.Phone == "" {
		return "", core_errors.ErrEmailAndPhoneEmpty
	}

	return input.Document, nil
}
