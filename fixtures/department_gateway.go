package fixtures

import (
	g "conformity-core/gateways/department_gateway"
)

type DepartmentGatewayMock struct {
	Create__Input            g.CreateDepartmentGatewayInput
	GetById__Output          *g.GetByIdDepartmentGatewayOutput
	GetUsersByFilter__Output *[]g.GetUsersByFilterGatewayOutput
}

func (d *DepartmentGatewayMock) Create(input g.CreateDepartmentGatewayInput) (string, error) {
	d.Create__Input = input
	return "id", nil
}

func (d *DepartmentGatewayMock) Delete(input g.DeleteInput) (bool, error) {
	return false, nil
}

func (d *DepartmentGatewayMock) GetOneByFilter(
	input g.GetOneByFilterInput,
) (*g.GetOneByFilterOutput, error) {
	return nil, nil
}

func (d *DepartmentGatewayMock) Patch(input g.PatchInput) (bool, error) {
	return false, nil
}

func (d *DepartmentGatewayMock) GetById(input string) (*g.GetByIdDepartmentGatewayOutput, error) {
	return d.GetById__Output, nil
}

func (d *DepartmentGatewayMock) GetUsersByFilter(input g.GetUsersByFilterGatewayInput) (*[]g.GetUsersByFilterGatewayOutput, error) {
	return d.GetUsersByFilter__Output, nil
}

func (d *DepartmentGatewayMock) GetUserByLogin(login string) (*g.UserDepartment, error) {
	return nil, nil
}

func (d *DepartmentGatewayMock) List(input g.ListInput) ([]g.ListOutput, error) {
	return nil, nil
}

func (u *DepartmentGatewayMock) PaginateUsers(input g.PaginateUsersInput) (*g.PaginateUsersOutput, error) {
	return nil, nil
}

func (u *DepartmentGatewayMock) CountMonthlyUsers(input g.CountMonthlyUsersInput) ([]g.CountMonthlyUsersOutput, error) {
	return nil, nil
}
