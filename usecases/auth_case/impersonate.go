package auth_case

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
)

type ImpersonateInput struct {
	Login string `validate:"required"`
}

func (a AuthCase) Impersonate(ctx *context.CoreCtx, input ImpersonateInput) (string, error) {
	if ctx.Session.Role != user_company_enum.RoleBackoffice {
		return "", core_errors.ErrWithoutPrivileges
	}

	user, err := a.departmentCase.GetUserByLogin(input.Login)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", nil
	}

	return generateToken(user, nil)
}
