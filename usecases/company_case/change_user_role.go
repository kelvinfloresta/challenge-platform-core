package company_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	g "conformity-core/gateways/company_gateway"
)

type ChangeUserRoleCaseInput struct {
	DepartmentID string                 `json:"department_id"`
	UserID       string                 `json:"user_id"`
	Role         user_company_enum.Role `json:"role"`
}

func (c CompanyCase) ChangeUserRole(ctx *context.CoreCtx, input ChangeUserRoleCaseInput) (updated bool, err error) {
	if ctx.Session.Role != user_company_enum.RoleCompanyManager &&
		ctx.Session.Role != user_company_enum.RoleBackoffice {
		return false, core_errors.ErrWithoutPrivileges
	}

	adapted := g.ChangeUserRoleInput{
		DepartmentID: input.DepartmentID,
		UserID:       input.UserID,
		Role:         input.Role,
	}

	updated, err = c.gateway.ChangeUserRole(adapted, nil)
	if err != nil {
		return
	}

	ctx.Logger.WithField("role", input.Role).Info("change-user-status")

	return updated, nil
}
