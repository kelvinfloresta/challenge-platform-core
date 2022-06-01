package department_case

import (
	g "conformity-core/gateways/department_gateway"
)

func (u DepartmentCase) GetById(id string) (*g.GetByIdDepartmentGatewayOutput, error) {
	return u.gateway.GetById(id)
}
