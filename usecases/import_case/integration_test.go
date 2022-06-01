package import_case_test

import (
	"conformity-core/enums/user_company_enum"
	"conformity-core/fixtures"
	"conformity-core/gateways/company_gateway"
	"conformity-core/gateways/import_gateway"
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/import_case"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ImportUsers__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	userCase := fixtures.UserCase.NewIntegration()
	gateway := import_gateway.GormImportGatewayFacade{DB: fixtures.DB_Test}
	sut := import_case.New(gateway, &userCase)

	companyId := fixtures.CreateCompany(t)
	err := sut.ImportUsers(fixtures.FakeBackofficeCtx, &import_gateway.ImportUsersInput{
		CompanyID: companyId,
		Schedule:  nil,
		Data: []import_gateway.ImportUsersData{
			{
				UserName:       "Test",
				Email:          "import@user.com",
				DepartmentName: "Importer",
			},
		},
	})

	require.Nil(t, err)
	departmentCase := fixtures.DepartmentCase.NewIntegration()
	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		CompanyID: companyId,
		Role:      user_company_enum.RoleCompanyManager,
	})

	page, err := departmentCase.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   5,
	})

	require.Nil(t, err)
	require.Len(t, page.Data, 1)
	require.Equal(t, page.Data[0].Name, "Test")
	require.Equal(t, page.Data[0].Email, "import@user.com")
	require.Equal(t, page.Data[0].DepartmentName, "Importer")
}

func Test_ImportUsers__Should_not_create_department_if_already_exists(t *testing.T) {
	fixtures.CleanTestDatabase()

	userCase := fixtures.UserCase.NewIntegration()
	gateway := import_gateway.GormImportGatewayFacade{DB: fixtures.DB_Test}
	sut := import_case.New(gateway, &userCase)

	companyId := fixtures.CreateCompany(t)
	departmentCase := fixtures.DepartmentCase.NewIntegration()

	departmentId, err := departmentCase.Create(fixtures.FakeBackofficeCtx, department_case.CreateDepartmentCaseInput{
		Name:      "Already Created",
		CompanyID: companyId,
	})
	require.Nil(t, err)

	err = sut.ImportUsers(fixtures.FakeBackofficeCtx, &import_gateway.ImportUsersInput{
		CompanyID: companyId,
		Schedule:  nil,
		Data: []import_gateway.ImportUsersData{
			{
				UserName:       "Test",
				Email:          "import@user.com",
				DepartmentName: "Already Created",
			},
		},
	})
	require.Nil(t, err)

	user := fixtures.CreateUser(t, &departmentId)
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		DepartmentID: user.DepartmentID,
		UserID:       user.ID,
	})

	departments, err := departmentCase.List(ctx)
	require.Nil(t, err)
	require.Len(t, departments, 1)
	require.Equal(t, "Already Created", departments[0].Name)
}
