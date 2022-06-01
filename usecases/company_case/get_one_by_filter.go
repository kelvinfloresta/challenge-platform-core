package company_case

import "conformity-core/gateways/company_gateway"

type GetOneByFilterInput struct {
	DepartmentID string
	Workspace    string
	Domain       string
}

func (c *CompanyCase) GetOneByFilter(input GetOneByFilterInput) (*company_gateway.GetOneByFilterOutput, error) {
	return c.gateway.GetOneByFilter(company_gateway.GetOneByFilterInput(input))
}
