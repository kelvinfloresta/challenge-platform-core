package company_case

import g "conformity-core/gateways/company_gateway"

func (u CompanyCase) GetUser(input g.GetUserInput) (*g.GetUserOutput, error) {
	return u.gateway.GetUser(input)
}
