package company_case

import g "conformity-core/gateways/company_gateway"

func (u CompanyCase) GetById(id string) (*g.GetByIdOutput, error) {
	return u.gateway.GetById(id)
}
