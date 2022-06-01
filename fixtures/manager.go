package fixtures

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	"conformity-core/gateways/company_gateway"
	"conformity-core/usecases/company_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type ChangeRoleToManagerInput struct {
	UserID       string
	DepartmentID string
}

func ChangeRoleToManager(t *testing.T, input ChangeRoleToManagerInput) *context.CoreCtx {
	companyCase := CompanyCase.NewIntegration()

	updated, err := companyCase.ChangeUserRole(FakeManagerCtx, company_case.ChangeUserRoleCaseInput{
		DepartmentID: input.DepartmentID,
		UserID:       input.UserID,
		Role:         user_company_enum.RoleCompanyManager,
	})

	require.Nil(t, err)
	require.True(t, updated)

	ctx := GetUserCtx(t, company_gateway.GetUserInput{
		DepartmentID: input.DepartmentID,
		UserID:       input.UserID,
	})

	return ctx
}
