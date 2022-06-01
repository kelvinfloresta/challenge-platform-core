package auth_case_test

import (
	core_errors "conformity-core/errors"
	"conformity-core/fixtures"
	"conformity-core/usecases/auth_case"
	"conformity-core/usecases/department_case"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoginSSO__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)

	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            user.OID,
		Email:          user.Email,
		Name:           "",
		DepartmentName: "",
		Workspace:      "",
	})

	require.Nil(t, err)
	require.NotEmpty(t, token)

	found, err := fixtures.DepartmentCase.NewIntegration().GetUsersByFilter(department_case.GetUsersByFilterCaseInput{})
	require.Nil(t, err)
	require.Len(t, *found, 1)
}

func Test_LoginSSO__Should_create_an_user_if_not_found(t *testing.T) {
	fixtures.CleanTestDatabase()

	company := fixtures.CreateCompanyV2(t)

	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          "Email",
		Name:           "",
		DepartmentName: "DepartmentName",
		Workspace:      company.Workspace,
	})

	require.Nil(t, err)
	require.NotEmpty(t, token)

	found, err := fixtures.DepartmentCase.NewIntegration().GetUsersByFilter(department_case.GetUsersByFilterCaseInput{})
	require.Nil(t, err)
	require.NotNil(t, found)
	require.Len(t, *found, 1)
}

func Test_LoginSSO__Should_create_a_department_if_not_found(t *testing.T) {
	fixtures.CleanTestDatabase()

	company := fixtures.CreateCompanyV2(t)

	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          "Email",
		Name:           "",
		DepartmentName: "DepartmentName",
		Workspace:      company.Workspace,
	})

	require.Nil(t, err)
	require.NotEmpty(t, token)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		CompanyID: company.ID,
	})
	found, err := fixtures.DepartmentCase.NewIntegration().List(ctx)

	require.Nil(t, err)
	require.Len(t, found, 1)
	require.Equal(t, "DepartmentName", found[0].Name)
}

func Test_LoginSSO__Should_not_create_a_department_if_found_one_with_the_same_name(t *testing.T) {
	fixtures.CleanTestDatabase()

	company := fixtures.CreateCompanyV2(t)
	departmentCase := fixtures.DepartmentCase.NewIntegration()
	_, err := departmentCase.Create(fixtures.FakeBackofficeCtx, department_case.CreateDepartmentCaseInput{
		CompanyID: company.ID,
		Name:      "DepartmentName",
	})
	require.Nil(t, err)

	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          "Email",
		Name:           "",
		DepartmentName: "DepartmentName",
		Workspace:      company.Workspace,
	})

	require.Nil(t, err)
	require.NotEmpty(t, token)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		CompanyID: company.ID,
	})

	found, err := departmentCase.List(ctx)

	require.Nil(t, err)
	require.Len(t, found, 1)
	require.Equal(t, "DepartmentName", found[0].Name)
}

func Test_LoginSSO__Should_search_department_name_in_case_insensitive(t *testing.T) {
	fixtures.CleanTestDatabase()

	company := fixtures.CreateCompanyV2(t)
	departmentCase := fixtures.DepartmentCase.NewIntegration()
	DEPARTMENT_WITH_DIFF_CASE := "departmentname"
	CORRECT_DEPARTMENT_NAME := "DepartmentName"

	_, err := departmentCase.Create(fixtures.FakeBackofficeCtx, department_case.CreateDepartmentCaseInput{
		CompanyID: company.ID,
		Name:      CORRECT_DEPARTMENT_NAME,
	})
	require.Nil(t, err)

	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          "Email",
		Name:           "",
		DepartmentName: DEPARTMENT_WITH_DIFF_CASE,
		Workspace:      company.Workspace,
	})

	require.Nil(t, err)
	require.NotEmpty(t, token)

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		CompanyID: company.ID,
	})

	found, err := departmentCase.List(ctx)

	require.Nil(t, err)
	require.Len(t, found, 1)
	require.Equal(t, CORRECT_DEPARTMENT_NAME, found[0].Name)
}

func Test_LoginSSO__Should_fail_if_email_is_empty(t *testing.T) {
	fixtures.CleanTestDatabase()

	EMPTY_EMAIL := ""
	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          EMPTY_EMAIL,
		Name:           "Name",
		DepartmentName: "DepartmentName",
		Workspace:      "",
	})

	require.EqualError(t, err, "MISSING_EMAIL")
	require.Empty(t, token)

}

func Test_LoginSSO__Should_fail_department_is_empty_and_is_a_new_user(t *testing.T) {
	fixtures.CleanTestDatabase()

	EMPTY_DEPARTMENT := ""
	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          "Email",
		Name:           "Name",
		DepartmentName: EMPTY_DEPARTMENT,
		Workspace:      "",
	})

	require.ErrorIs(t, err, core_errors.ErrUserWithoutDepartment)
	require.Empty(t, token)
}

func Test_LoginSSO__Should_fail_if_company_not_found(t *testing.T) {
	fixtures.CleanTestDatabase()

	NON_EXISTENT_COMPANY := fixtures.UUID(t)
	sut := fixtures.AuthCase.NewIntegration()
	token, err := sut.LoginSSO(auth_case.LoginSSOInput{
		OID:            "OID",
		Email:          "Email",
		Name:           "Name",
		DepartmentName: "DepartmentName",
		Workspace:      NON_EXISTENT_COMPANY,
	})

	require.ErrorIs(t, err, core_errors.ErrCompanyNotFound)
	require.Empty(t, token)
}
