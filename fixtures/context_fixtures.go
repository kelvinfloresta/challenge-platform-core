package fixtures

import (
	core "conformity-core/context"
	"conformity-core/enums/user_company_enum"
	"conformity-core/gateways/company_gateway"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type CreateContextInput struct {
	UserID       string
	DepartmentID string
	CompanyID    string
	Role         user_company_enum.Role
}

func GetUserCtx(t *testing.T, input company_gateway.GetUserInput) *core.CoreCtx {
	companyCase := CompanyCase.NewIntegration()
	user, err := companyCase.GetUser(input)

	require.Nil(t, err)

	departmentCase := DepartmentCase.NewIntegration()
	department, err := departmentCase.GetById(input.DepartmentID)

	require.Nil(t, err)
	require.NotNil(t, department)

	return core.New(context.WithValue(context.Background(), core.SessionKey,
		&core.UserSession{
			UserID:       user.UserID,
			DepartmentID: user.DepartmentID,
			CompanyID:    department.CompanyID,
			Role:         user.Role,
		}))
}

func CreateContext(input CreateContextInput) *core.CoreCtx {
	return core.New(context.WithValue(context.Background(), core.SessionKey,
		&core.UserSession{
			UserID:       input.UserID,
			DepartmentID: input.DepartmentID,
			CompanyID:    input.CompanyID,
			Name:         "any_name",
			Login:        "any_login",
			Role:         input.Role,
		}))
}

var FakeManagerCtx = core.New(context.WithValue(context.Background(), core.SessionKey,
	&core.UserSession{
		Role: user_company_enum.RoleCompanyManager,
	}))

var FakeBackofficeCtx = core.New(context.WithValue(context.Background(), core.SessionKey,
	&core.UserSession{
		Role: user_company_enum.RoleBackoffice,
	}))
