package fixtures

import (
	"conformity-core/gateways/department_gateway"
	"conformity-core/usecases/department_case"
	"testing"

	"github.com/stretchr/testify/require"
)

type departmentCaseFixture struct{}

var DepartmentCase departmentCaseFixture = departmentCaseFixture{}

func (u departmentCaseFixture) NewUnit() (department_case.DepartmentCase, *DepartmentGatewayMock) {
	gatewayMock := &DepartmentGatewayMock{}
	departmentCase := department_case.New(gatewayMock)

	return departmentCase, gatewayMock
}

func (u departmentCaseFixture) NewIntegration() department_case.DepartmentCase {
	gateway := &department_gateway.GormDepartmentGatewayFacade{DB: DB_Test}
	return department_case.New(gateway)
}

func CreateDepartment(t *testing.T, companyId string) string {
	departmentCase := DepartmentCase.NewIntegration()

	id, err := departmentCase.Create(FakeBackofficeCtx, department_case.CreateDepartmentCaseInput{
		Name:      "Any name",
		CompanyID: companyId,
	})

	require.Nil(t, err)

	return id
}
