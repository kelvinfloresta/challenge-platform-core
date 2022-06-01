package fixtures

import (
	"conformity-core/context"
	g "conformity-core/gateways/company_gateway"
	"conformity-core/usecases/company_case"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type companyCaseFixture struct{}

var CompanyCase companyCaseFixture = companyCaseFixture{}

func (u companyCaseFixture) NewUnit() (company_case.CompanyCase, *CompanyGatewayMock) {
	gatewayMock := &CompanyGatewayMock{}
	companyCase := company_case.New(gatewayMock)

	return companyCase, gatewayMock
}

func (u companyCaseFixture) NewIntegration() company_case.CompanyCase {
	gateway := &g.GormCompanyGatewayFacade{DB: DB_Test}
	return company_case.New(gateway)
}

type CompanyGatewayMock struct {
	Create__Input          g.CreateCompanyGatewayInput
	GetOneByFilter__Output *g.GetOneByFilterOutput
	GetById__Output        *g.GetByIdOutput
}

func (u *CompanyGatewayMock) Create(input g.CreateCompanyGatewayInput) (string, error) {
	u.Create__Input = input
	return "id", nil
}

func (u *CompanyGatewayMock) GetById(input string) (*g.GetByIdOutput, error) {
	return u.GetById__Output, nil
}

func (c *CompanyGatewayMock) GetOneByFilter(input g.GetOneByFilterInput) (*g.GetOneByFilterOutput, error) {
	return c.GetOneByFilter__Output, nil
}

func (u *CompanyGatewayMock) GetUser(input g.GetUserInput) (*g.GetUserOutput, error) {
	return nil, nil
}

func (u *CompanyGatewayMock) GetTransaction(ctx *context.CoreCtx) *gorm.DB {
	return nil
}

func (u *CompanyGatewayMock) ChangeUserRole(input g.ChangeUserRoleInput, tx *gorm.DB) (bool, error) {
	return false, nil
}

func (c *CompanyGatewayMock) ListManagers(input g.ListManagersInput) ([]g.ListManagersOutput, error) {
	return nil, nil
}

func CreateCompany(t *testing.T) string {
	companyCase := CompanyCase.NewIntegration()

	id, err := companyCase.Create(DUMMY_CONTEXT, company_case.CreateCompanyCaseInput{
		Name:      "Any name",
		Document:  "Any document",
		Workspace: UUID(t),
		Domain:    UUID(t),
	})

	require.Nil(t, err)

	return id
}

type CreateCompanyOutput struct {
	ID        string
	Workspace string
}

func CreateCompanyV2(t *testing.T) (out CreateCompanyOutput) {
	companyCase := CompanyCase.NewIntegration()

	workspace := UUID(t)
	id, err := companyCase.Create(DUMMY_CONTEXT, company_case.CreateCompanyCaseInput{
		Name:      "Any name",
		Document:  "Any document",
		Workspace: workspace,
	})
	require.Nil(t, err)

	out.Workspace = workspace
	out.ID = id

	return
}
