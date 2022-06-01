package department_case

import g "conformity-core/gateways/department_gateway"

type GetOneByFilter struct {
	Name      string
	CompanyID string
}

func (d DepartmentCase) GetOneByFilter(input GetOneByFilter) (*g.GetOneByFilterOutput, error) {
	return d.gateway.GetOneByFilter(g.GetOneByFilterInput(input))
}
