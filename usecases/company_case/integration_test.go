package company_case_test

import (
	"conformity-core/enums/user_company_enum"
	"conformity-core/fixtures"
	"conformity-core/gateways/company_gateway"
	"conformity-core/usecases/company_case"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompanyCase_Create__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CompanyCase.NewIntegration()
	companyToCreate := company_case.CreateCompanyCaseInput{
		Name:            "Name",
		Document:        "Document",
		Workspace:       "Workspace",
		RequirePassword: true,
	}

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, companyToCreate)
	require.Nil(t, err)

	found, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Equal(t, found.ID, id)
	assert.Equal(t, found.Name, companyToCreate.Name)
	assert.Equal(t, found.Document, companyToCreate.Document)
	assert.Equal(t, found.Workspace, companyToCreate.Workspace)
	assert.Equal(t, found.RequirePassword, companyToCreate.RequirePassword)
}

func TestCompanyCase_Create__Should_define_CreatedAt(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CompanyCase.NewIntegration()
	companyToCreate := company_case.CreateCompanyCaseInput{
		Name:     "Name",
		Document: "Document",
	}

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, companyToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.IsType(t, time.Time{}, result.CreatedAt)
}

func TestCompanyCase_Create__Should_define_UpdatedAt(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CompanyCase.NewIntegration()
	companyToCreate := company_case.CreateCompanyCaseInput{
		Name:     "Name",
		Document: "Document",
	}

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, companyToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.IsType(t, time.Time{}, result.UpdatedAt)
}

func TestCompanyCase_Create__Should_define_an_empty_DeletedAt(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CompanyCase.NewIntegration()
	companyToCreate := company_case.CreateCompanyCaseInput{
		Name:     "Name",
		Document: "Document",
	}

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, companyToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Empty(t, result.DeletedAt)
}

func TestCompanyCase_Create__Should_require_password_be_optional(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.CompanyCase.NewIntegration()
	companyToCreate := company_case.CreateCompanyCaseInput{
		Name:     "Name",
		Document: "Document",
	}

	id, err := sut.Create(fixtures.DUMMY_CONTEXT, companyToCreate)
	require.Nil(t, err)

	found, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Equal(t, found.RequirePassword, false)
}

func Test_ChangeUserRole_Should_do_happy_path(t *testing.T) {
	t.SkipNow()
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)

	sut := fixtures.CompanyCase.NewIntegration()

	updated, err := sut.ChangeUserRole(fixtures.DUMMY_CONTEXT, company_case.ChangeUserRoleCaseInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
		Role:         user_company_enum.RoleCompanyManager,
	})

	require.Nil(t, err)
	require.True(t, updated)

	companyUser, err := sut.GetUser(company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.Equal(t, user_company_enum.RoleCompanyManager, companyUser.Role)
}

func Test_GetOneByFilter_Should_filter_by_department(t *testing.T) {
	fixtures.CleanTestDatabase()

	fixtures.CreateUser(t, nil)
	user := fixtures.CreateUser(t, nil)

	sut := fixtures.CompanyCase.NewIntegration()
	companyFound, err := sut.GetOneByFilter(company_case.GetOneByFilterInput{
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.Equal(t, user.CompanyID, companyFound.ID)
}
