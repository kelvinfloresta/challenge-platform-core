package user_case_test

import (
	"conformity-core/enums/campaign_enum"
	core_errors "conformity-core/errors"
	"conformity-core/fixtures"
	"conformity-core/gateways/campaign_gateway"
	"conformity-core/usecases/user_case"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Create__Should_do_happy_path(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.UserCase.NewIntegration()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	userToCreate := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        "Email",
		Login:        "Login",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	id, err := sut.Create(fixtures.FakeManagerCtx, userToCreate)
	require.Nil(t, err)

	found, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Equal(t, found.ID, id)
	assert.Equal(t, found.Name, userToCreate.Name)
	assert.Equal(t, found.Email, userToCreate.Email)
	assert.Equal(t, found.Document, userToCreate.Document)
	assert.Equal(t, found.Phone, userToCreate.Phone)
	assert.Equal(t, found.BirthDate, userToCreate.BirthDate)
	assert.Equal(t, found.JobPosition, userToCreate.JobPosition)
}

func Test_Create__Should_fail_if_email_already_exists(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.UserCase.NewIntegration()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	DUPLICATED_EMAIL := "Email"
	user_1 := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        DUPLICATED_EMAIL,
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	user_2 := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        DUPLICATED_EMAIL,
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	sut.Create(fixtures.FakeManagerCtx, user_1)
	_, err := sut.Create(fixtures.FakeManagerCtx, user_2)
	require.ErrorIs(t, core_errors.ErrDuplicatedEmail, err, "Should fail because this email already exists")
}

func Test_Create__Should_not_fail_if_both_email_is_empty(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	sut := fixtures.UserCase.NewIntegration()

	EMPTY_EMAIL := ""
	user_1 := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        EMPTY_EMAIL,
		Login:        "login2",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	user_2 := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        EMPTY_EMAIL,
		Login:        "login1",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	_, err := sut.Create(fixtures.FakeManagerCtx, user_1)
	require.Nil(t, err)
	_, err = sut.Create(fixtures.FakeManagerCtx, user_2)
	require.Nil(t, err)
}

func Test_Create__Should_use_email_as_login(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	sut := fixtures.UserCase.NewIntegration()

	_, err := sut.Create(fixtures.FakeManagerCtx, user_case.CreateInput{
		Name:         "Name",
		Email:        "Email",
		DepartmentID: departmentId,
	})
	require.Nil(t, err)

	user, err := fixtures.DepartmentCase.NewIntegration().GetUserByLogin("Email")
	require.Nil(t, err)

	require.NotNil(t, user)
	require.Equal(t, user.Login, "Email")
}

func Test_Create__Should_use_document_as_login_if_email_is_empty(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	sut := fixtures.UserCase.NewIntegration()

	_, err := sut.Create(fixtures.FakeManagerCtx, user_case.CreateInput{
		Name:         "Name",
		Email:        "",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	})
	require.Nil(t, err)

	user, err := fixtures.DepartmentCase.NewIntegration().GetUserByLogin("Document")
	require.Nil(t, err)

	require.NotNil(t, user)
	require.Equal(t, user.Login, "Document")
}

func Test_Create__Should_match_password(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	sut := fixtures.UserCase.NewIntegration()
	userToCreate := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        "Email",
		Login:        "",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	id, err := sut.Create(fixtures.FakeManagerCtx, userToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	fixtures.MatchPassword(t, result.Password, userToCreate.Password)
}

func Test_Create__Should_define_CreatedAt(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)

	sut := fixtures.UserCase.NewIntegration()
	userToCreate := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        "Email",
		Login:        "",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	id, err := sut.Create(fixtures.FakeManagerCtx, userToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.IsType(t, time.Time{}, result.CreatedAt)
}

func Test_Create__Should_define_UpdatedAt(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)

	sut := fixtures.UserCase.NewIntegration()
	userToCreate := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        "Email",
		Login:        "",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	id, err := sut.Create(fixtures.FakeManagerCtx, userToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.IsType(t, time.Time{}, result.UpdatedAt)
}

func Test_Create__Should_define_an_empty_DeletedAt(t *testing.T) {
	fixtures.CleanTestDatabase()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)

	sut := fixtures.UserCase.NewIntegration()
	userToCreate := user_case.CreateInput{
		Name:         "Name",
		Password:     "12345678",
		Email:        "Email",
		Login:        "",
		DepartmentID: departmentId,
		Document:     "Document",
		Phone:        "Phone",
		BirthDate:    "BirthDate",
		JobPosition:  "JobPosition",
	}

	id, err := sut.Create(fixtures.FakeManagerCtx, userToCreate)
	require.Nil(t, err)

	result, err := sut.GetById(id)
	require.Nil(t, err)

	assert.Empty(t, result.DeletedAt)
}

func Test_Create__Should_not_create_if_department_does_not_exist(t *testing.T) {
	fixtures.CleanTestDatabase()

	sut := fixtures.UserCase.NewIntegration()
	userId, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Name:         "Rabi",
		Email:        "rabi@teo.gov",
		DepartmentID: fixtures.UUID(t),
	})
	require.NotNil(t, err)

	result, err := sut.GetById(userId)

	require.Nil(t, err)
	require.Nil(t, result)
}

func Test_Create__Should_include_user_in_the_campaign(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)
	sut := fixtures.UserCase.NewIntegration()
	newUserId, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Name:         "New User",
		Email:        "new@user.com",
		DepartmentID: campaign.DepartmentID,
	})
	require.Nil(t, err)

	campaignCase := fixtures.CampaignCase.NewIntegration()
	result, err := campaignCase.GetResult(fixtures.FakeBackofficeCtx, campaign_gateway.GetResultInput{
		UserID:       newUserId,
		DepartmentID: campaign.DepartmentID,
		ChallengeID:  campaign.Challenge.ID,
		CampaignID:   campaign.CampaignID,
		QuestionID:   campaign.Challenge.Questions[0].ID,
	})
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, campaign.CampaignID, result.CampaignID)
	require.Equal(t, campaign.Challenge.ID, result.ChallengeID)
	require.False(t, result.Correct)
	require.Equal(t, campaign.DepartmentID, result.DepartmentID)
	require.Equal(t, campaign.Challenge.Questions[0].ID, result.QuestionID)
	require.Equal(t, campaign_enum.Active, result.Status)
	require.Zero(t, result.Tries)
	require.Equal(t, newUserId, result.UserID)
}

func Test_Create__Should_not_include_user_in_campaigns_from_another_company(t *testing.T) {
	fixtures.CleanTestDatabase()

	campaign := fixtures.CreateDefaultCampaign(t, nil)

	sut := fixtures.UserCase.NewIntegration()
	companyId := fixtures.CreateCompany(t)
	departmentId := fixtures.CreateDepartment(t, companyId)
	DIFF_USER_COMPANY, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Name:         "Diff Company",
		Email:        "diff-company@user.com",
		DepartmentID: departmentId,
	})
	require.Nil(t, err)

	campaignCase := fixtures.CampaignCase.NewIntegration()
	result, err := campaignCase.GetResult(fixtures.FakeBackofficeCtx, campaign_gateway.GetResultInput{
		UserID:       DIFF_USER_COMPANY,
		DepartmentID: campaign.DepartmentID,
		ChallengeID:  campaign.Challenge.ID,
		CampaignID:   campaign.CampaignID,
		QuestionID:   campaign.Challenge.Questions[0].ID,
	})
	require.Nil(t, err)
	require.Nil(t, result)
}
