package department_case_test

import (
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/fixtures"
	"conformity-core/gateways/company_gateway"
	"conformity-core/usecases/campaign_case"
	"conformity-core/usecases/department_case"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepartmentCase_Create__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)

	sut := fixtures.DepartmentCase.NewIntegration()
	departmentToCreate := department_case.CreateDepartmentCaseInput{
		Name:      "Name",
		CompanyID: companyId,
	}

	id, err := sut.Create(fixtures.FakeBackofficeCtx, departmentToCreate)
	require.Nil(t, err)

	found, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Equal(t, found.ID, id)
	assert.Equal(t, found.Name, departmentToCreate.Name)
}

func TestDepartmentCase_Create__Should_define_CreatedAt(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	sut := fixtures.DepartmentCase.NewIntegration()
	departmentToCreate := department_case.CreateDepartmentCaseInput{
		Name:      "Name",
		CompanyID: companyId,
	}

	id, err := sut.Create(fixtures.FakeBackofficeCtx, departmentToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.IsType(t, time.Time{}, result.CreatedAt)
}

func TestDepartmentCase_Create__Should_define_UpdatedAt(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	sut := fixtures.DepartmentCase.NewIntegration()
	departmentToCreate := department_case.CreateDepartmentCaseInput{
		Name:      "Name",
		CompanyID: companyId,
	}

	id, err := sut.Create(fixtures.FakeBackofficeCtx, departmentToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.IsType(t, time.Time{}, result.UpdatedAt)
}

func TestDepartmentCase_Create__Should_define_an_empty_DeletedAt(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	sut := fixtures.DepartmentCase.NewIntegration()
	departmentToCreate := department_case.CreateDepartmentCaseInput{
		Name:      "Name",
		CompanyID: companyId,
	}

	id, err := sut.Create(fixtures.FakeBackofficeCtx, departmentToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Empty(t, result.DeletedAt)
}

func Test_Delete__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	anotherDepartment := fixtures.CreateDepartment(t, user.CompanyID)
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	sut := fixtures.DepartmentCase.NewIntegration()
	deleted, err := sut.Delete(ctx, anotherDepartment)
	require.Nil(t, err)
	require.True(t, deleted)

	department, err := sut.GetById(anotherDepartment)
	require.Nil(t, err)
	require.Nil(t, department)
}

func Test_Delete__Should_not_delete_if_have_users_in_department(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	sut := fixtures.DepartmentCase.NewIntegration()
	deleted, err := sut.Delete(ctx, user.DepartmentID)
	require.ErrorIs(t, err, core_errors.ErrDeleteDepartmentWithUsers)
	require.False(t, deleted)

	department, err := sut.GetById(user.DepartmentID)
	require.Nil(t, err)
	require.NotNil(t, department)
}

func Test_Delete__Should_not_delete_if_not_found(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	departmentId := fixtures.UUID(t)
	sut := fixtures.DepartmentCase.NewIntegration()
	deleted, err := sut.Delete(ctx, departmentId)
	require.Nil(t, err)
	require.False(t, deleted)
}

func Test_Delete__Should_not_delete_if_user_does_not_belong_to_the_company(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	OTHER_COMPANYS_USER := fixtures.CreateUser(t, nil)

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       OTHER_COMPANYS_USER.ID,
		DepartmentID: OTHER_COMPANYS_USER.DepartmentID,
	})

	sut := fixtures.DepartmentCase.NewIntegration()
	deleted, err := sut.Delete(ctx, departmentId)
	require.Nil(t, err)
	require.False(t, deleted)

	department, err := sut.GetById(departmentId)
	require.Nil(t, err)
	require.NotNil(t, department)
}

func Test_PaginateUsers_Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)

	sut := fixtures.DepartmentCase.NewIntegration()

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
		CompanyID:    user.CompanyID,
	})

	users, err := sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   10,
	})

	require.Nil(t, err)
	require.Equal(t, 1, users.MaxPages)
	require.NotEmpty(t, users.Data)

	for _, user := range users.Data {
		require.NotEmpty(t, user.ID)
		require.Equal(t, "Any name", user.Name)
		require.Equal(t, "Any name", user.DepartmentName)
		require.Equal(t, user_company_enum.Active, user.Status)
	}
}

func Test_PaginateUsers_Should_not_list_deleted_users(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)

	sut := fixtures.DepartmentCase.NewIntegration()

	ctx := fixtures.ChangeRoleToManager(t, fixtures.ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	removed, err := fixtures.CampaignCase.NewIntegration().RemoveUser(fixtures.FakeBackofficeCtx, campaign_case.RemoveUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	require.Nil(t, err)
	require.True(t, removed)

	users, err := sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   10,
	})

	require.Nil(t, err)
	require.Equal(t, 0, users.MaxPages)
	require.Empty(t, users.Data)

}

func Test_PaginateUsers_Should_paginate(t *testing.T) {
	fixtures.CleanTestDatabase()

	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)

	for i := 0; i < 5; i++ {
		fixtures.CreateUser(t, &departmentId)
	}

	sut := fixtures.DepartmentCase.NewIntegration()
	expectedPages := 3

	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		DepartmentID: departmentId,
		CompanyID:    companyId,
	})

	allUsers, err := sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   5,
	})

	require.Nil(t, err)
	require.Len(t, allUsers.Data, 5)

	count := 0
	for actualPage := 0; actualPage < expectedPages; actualPage++ {
		users, err := sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
			ActualPage: actualPage,
			PageSize:   2,
		})

		require.Nil(t, err)
		for _, user := range users.Data {
			expectedUser := (allUsers.Data)[count]
			require.Equal(t, expectedUser, user)
			count++
		}

		require.Equal(t, expectedPages, users.MaxPages)
		isLastPage := actualPage == 2
		if isLastPage {
			require.Len(t, users.Data, 1)
		}

		if !isLastPage {
			require.Len(t, users.Data, 2)
		}
	}
}

func Test_PaginateUsers_Should_filter_by_name(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	sut := fixtures.DepartmentCase.NewIntegration()
	result, err := sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   5,
		Name:       "non existent",
	})

	require.Nil(t, err)
	require.Len(t, result.Data, 0)

	result, err = sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   5,
		Name:       "ny",
	})

	require.Nil(t, err)
	require.Len(t, result.Data, 1)
	require.Equal(t, result.Data[0].Name, "Any name")
}

func Test_PaginateUsers_Should_filter_by_email(t *testing.T) {
	fixtures.CleanTestDatabase()

	user := fixtures.CreateUser(t, nil)
	ctx := fixtures.GetUserCtx(t, company_gateway.GetUserInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	sut := fixtures.DepartmentCase.NewIntegration()
	result, err := sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   5,
		Email:      "non existent",
	})

	require.Nil(t, err)
	require.Len(t, result.Data, 0)

	result, err = sut.PaginateUsers(ctx, department_case.PaginateUsersInput{
		ActualPage: 0,
		PageSize:   5,
		Email:      user.Email[3:6],
	})

	require.Nil(t, err)
	require.Len(t, result.Data, 1)
	require.Equal(t, result.Data[0].Name, "Any name")
}
