package department_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	g "conformity-core/gateways/department_gateway"

	"github.com/sirupsen/logrus"
)

type CreateDepartmentCaseInput struct {
	Name      string `json:"name" validate:"required"`
	CompanyID string `json:"company_id"`
}

func (u DepartmentCase) Create(ctx *context.CoreCtx, input CreateDepartmentCaseInput) (string, error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager && ctx.Session.Role != user_company_enum.RoleBackoffice {
		return "", core_errors.ErrWithoutPrivileges
	}

	if ctx.Session.Role != user_company_enum.RoleBackoffice || input.CompanyID == "" {
		input.CompanyID = ctx.Session.CompanyID
	}

	id, err := u.gateway.Create(g.CreateDepartmentGatewayInput(input))

	if err != nil {
		return "", err
	}

	ctx.Logger.WithFields(logrus.Fields{
		"id":        id,
		"name":      input.Name,
		"companyId": input.CompanyID,
	}).Info("department-created")

	return id, nil
}
