package company_case

import (
	"conformity-core/gateways/company_gateway"
)

type CompanyCase struct {
	gateway company_gateway.ICompanyGateway
}

func New(gateway company_gateway.ICompanyGateway) CompanyCase {
	return CompanyCase{gateway}
}
