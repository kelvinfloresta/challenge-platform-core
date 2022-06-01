package user_case_test

import (
	"conformity-core/config"
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/fixtures"
	"conformity-core/gateways/company_gateway"
	"conformity-core/usecases/user_case"
	"conformity-core/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCase_Create__Should_encrypt_password(t *testing.T) {
	userCase, gatewayMock := fixtures.UserCase.NewUnit()
	const MY_PASSWORD = "12345678"

	_, err := userCase.Create(fixtures.FakeManagerCtx, user_case.CreateInput{
		Password: MY_PASSWORD,
		Email:    "email",
	})

	require.Nil(t, err)

	encryptedPassword := gatewayMock.Create__Input.Password
	assert.NotEqualf(t, MY_PASSWORD, encryptedPassword, "Password: %s; Want != %s", encryptedPassword, MY_PASSWORD)
}

func TestUserCase_Create__Should_be_able_to_decrypt_the_password(t *testing.T) {
	userCase, gatewayMock := fixtures.UserCase.NewUnit()
	const MY_PASSWORD = "12345678"
	_, err := userCase.Create(fixtures.FakeManagerCtx, user_case.CreateInput{
		Password: MY_PASSWORD,
		Email:    "email",
	})

	require.Nil(t, err)

	encryptedPassword := gatewayMock.Create__Input.Password
	descryptError := utils.ComparePassword(MY_PASSWORD, encryptedPassword)

	assert.Nilf(t, descryptError, "Failed to decrypt password %s", encryptedPassword)
}

func Test_Create__Should_fail_if_role_is_user(t *testing.T) {
	sut, _ := fixtures.UserCase.NewUnit()
	ctx := fixtures.CreateContext(fixtures.CreateContextInput{
		UserID:       "UserID",
		DepartmentID: "DepartmentID",
		CompanyID:    "CompanyID",
		Role:         user_company_enum.RoleUser,
	})
	_, err := sut.Create(ctx, user_case.CreateInput{})
	require.ErrorIs(t, err, core_errors.ErrWithoutPrivileges)
}

func Test_Create__Should_allow_companyManager_role(t *testing.T) {
	sut, _ := fixtures.UserCase.NewUnit()
	_, err := sut.Create(fixtures.FakeManagerCtx, user_case.CreateInput{
		Email: "email",
	})
	require.Nil(t, err)
}

func Test_Create__Should_allow_backoffice_role(t *testing.T) {
	sut, _ := fixtures.UserCase.NewUnit()
	_, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Email: "email",
	})
	require.Nil(t, err)
}

func Test_Create__Should_fail_if_email_and_document_is_empty(t *testing.T) {
	sut, _ := fixtures.UserCase.NewUnit()
	_, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Email:    "",
		Document: "",
	})

	require.ErrorIs(t, err, core_errors.ErrEmailAndDocumentEmpty)
}

func Test_Create__Should_fail_if_email_and_phone_is_empty(t *testing.T) {
	sut, _ := fixtures.UserCase.NewUnit()
	_, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Email:    "",
		Document: "Document",
	})

	require.ErrorIs(t, err, core_errors.ErrEmailAndPhoneEmpty)
}

func Test_Create__Should_not_fail_if_fill_document_and_phone_but_not_email(t *testing.T) {
	sut, _ := fixtures.UserCase.NewUnit()
	_, err := sut.Create(fixtures.FakeBackofficeCtx, user_case.CreateInput{
		Email:    "",
		Document: "Document",
		Phone:    "Phone",
	})

	require.Nil(t, err)
}

func Test_NotifyNewAccount__Should_send_correctly_template_if_company_requires_password(t *testing.T) {
	gatewayMock := &fixtures.UserGatewayMock{}
	notificationCase, notificationGateway := fixtures.NotificationCase.NewUnit()
	companyCase, companyGateway := fixtures.CompanyCase.NewUnit()
	campaignCase, _ := fixtures.CampaignCase.NewUnit()
	companyGateway.GetOneByFilter__Output = &company_gateway.GetOneByFilterOutput{
		RequirePassword: true,
	}

	sut := user_case.New(gatewayMock, notificationCase, &companyCase, &campaignCase)

	notification := user_case.NotifyNewAccountInput{
		Name:     "Name",
		Email:    "Email",
		Login:    "Login",
		Schedule: nil,
	}

	sut.NotifyNewAccount(&notification)

	require.Len(t, notificationGateway.Messages, 1)

	for _, message := range notificationGateway.Messages {
		require.Equal(t, "new-account", message.TemplateID)
		require.Equal(t, []string{notification.Email}, message.To)
		require.Equal(t, "Seja bem-vindo ao Conformity Pro", message.Subject)
		require.Equal(t, notification.Name, message.Variables["name"])
		require.Equal(t, notification.Login, message.Variables["login"])
		require.Equal(t, config.AppURL, message.Variables["link"])
		require.Nil(t, message.Schedule)
	}
}

func Test_NotifyNewAccount__Should_send_correctly_template_if_company_does_not_require_password(t *testing.T) {
	gatewayMock := &fixtures.UserGatewayMock{}
	notificationCase, notificationGateway := fixtures.NotificationCase.NewUnit()
	companyCase, companyGateway := fixtures.CompanyCase.NewUnit()
	companyGateway.GetOneByFilter__Output = &company_gateway.GetOneByFilterOutput{
		RequirePassword: false,
	}
	campaignCase, _ := fixtures.CampaignCase.NewUnit()

	sut := user_case.New(gatewayMock, notificationCase, &companyCase, &campaignCase)

	notification := user_case.NotifyNewAccountInput{
		Name:     "Name",
		Email:    "Email",
		Login:    "Login",
		Schedule: nil,
	}

	sut.NotifyNewAccount(&notification)

	require.Len(t, notificationGateway.Messages, 1)

	for _, message := range notificationGateway.Messages {
		require.Equal(t, "new-account-without-password", message.TemplateID)
		require.Equal(t, []string{notification.Email}, message.To)
		require.Equal(t, "Seja bem-vindo ao Conformity Pro", message.Subject)
		require.Equal(t, notification.Name, message.Variables["name"])
		require.Equal(t, notification.Login, message.Variables["login"])
		require.Equal(t, config.AppURL, message.Variables["link"])
		require.Nil(t, message.Schedule)
	}
}
