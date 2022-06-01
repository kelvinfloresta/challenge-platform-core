package company_case

import (
	"conformity-core/enums/user_company_enum"
	g "conformity-core/gateways/company_gateway"
)

type ListManagersInput struct {
	CompanyIDs []string
}

func (u CompanyCase) ListManagers(companyIds []string) ([]g.ListManagersOutput, error) {
	return u.gateway.ListManagers(g.ListManagersInput{
		CompanyIDs: companyIds,
		Status:     user_company_enum.Active,
	})
}
