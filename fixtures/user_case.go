package fixtures

import (
	"conformity-core/context"
	g "conformity-core/gateways/user_gateway"
	"conformity-core/usecases/user_case"
	"conformity-core/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

type userCaseFixture struct{}

var UserCase userCaseFixture = userCaseFixture{}

func (u userCaseFixture) NewUnit() (user_case.UserCase, *UserGatewayMock) {
	gatewayMock := &UserGatewayMock{}
	notificationCase, _ := NotificationCase.NewUnit()
	companyCase, _ := CompanyCase.NewUnit()
	campaignCase, _ := CampaignCase.NewUnit()
	userCase := user_case.New(gatewayMock, notificationCase, &companyCase, &campaignCase)
	return userCase, gatewayMock
}

func (u userCaseFixture) NewIntegration() user_case.UserCase {
	gateway := &g.GormUserGatewayFacade{DB: DB_Test}
	notificationCase, _ := NotificationCase.NewUnit()
	companyCase := CompanyCase.NewIntegration()
	campaignCase := CampaignCase.NewIntegration()
	return user_case.New(gateway, notificationCase, &companyCase, &campaignCase)
}

func MatchPassword(t *testing.T, hashedPassword string, password string) {
	err := utils.ComparePassword(password, hashedPassword)
	require.Nil(t, err)

}

type UserGatewayMock struct {
	Create__Input      g.CreateInput
	GetById__Output    *g.GetByIdUserGatewayOutput
	GetByEmail__Output *g.GetByEmailUserGatewayOutput
}

func (u *UserGatewayMock) Create(input g.CreateInput) (string, error) {
	u.Create__Input = input
	return "id", nil
}

func (u *UserGatewayMock) Patch(input g.PatchInput) (bool, error) {
	return false, nil
}

func (u *UserGatewayMock) GetOneByFilter(input g.GetOneByFilterInput) (*g.GetOneByFilterOutput, error) {
	return nil, nil
}

func (u *UserGatewayMock) GetById(input string) (*g.GetByIdUserGatewayOutput, error) {
	return u.GetById__Output, nil
}

func (u *UserGatewayMock) GetByEmail(input string) (*g.GetByEmailUserGatewayOutput, error) {
	return u.GetByEmail__Output, nil
}

func (u *UserGatewayMock) ChangePassword(input g.ChangePasswordInput) (bool, error) {
	return false, nil
}

type CreateUserOutput struct {
	ID           string
	OID          string
	Email        string
	DepartmentID string
	CompanyID    string
}

func CreateUser(t *testing.T, departmentId *string) (out CreateUserOutput) {
	userCase := UserCase.NewIntegration()

	if departmentId == nil {
		out.CompanyID = CreateCompany(t)
		out.DepartmentID = CreateDepartment(t, out.CompanyID)
	} else {
		out.DepartmentID = *departmentId
	}

	oid := UUID(t)
	email := UUID(t)
	id, err := userCase.Create(FakeManagerCtx, user_case.CreateInput{
		Name:         "Any name",
		Email:        email,
		OID:          oid,
		DepartmentID: out.DepartmentID,
	})

	require.Nil(t, err)

	out.ID = id
	out.OID = oid
	out.Email = email

	return
}

type CreateCompanyManagerOutput struct {
	ID           string
	OID          string
	Email        string
	DepartmentID string
	CompanyID    string
	Ctx          *context.CoreCtx
}

func CreateCompanyManager(t *testing.T, departmentId *string) (out CreateCompanyManagerOutput) {
	user := CreateUser(t, departmentId)

	ctx := ChangeRoleToManager(t, ChangeRoleToManagerInput{
		UserID:       user.ID,
		DepartmentID: user.DepartmentID,
	})

	out.ID = user.ID
	out.OID = user.OID
	out.Email = user.Email
	out.DepartmentID = user.DepartmentID
	out.CompanyID = user.CompanyID
	out.Ctx = ctx

	return
}
