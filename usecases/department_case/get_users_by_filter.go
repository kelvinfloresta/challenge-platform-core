package department_case

import (
	"conformity-core/enums/user_company_enum"
	g "conformity-core/gateways/department_gateway"
	"time"
)

type GetUsersByFilterCaseInput struct {
	CompanyID         string
	Departments       []string
	UserCreatedAtGTE  time.Time
	UserCreatedAtLTE  time.Time
	UserCompanyStatus user_company_enum.Status
}

func (u DepartmentCase) GetUsersByFilter(input GetUsersByFilterCaseInput) (*[]g.GetUsersByFilterGatewayOutput, error) {
	return u.gateway.GetUsersByFilter(g.GetUsersByFilterGatewayInput(input))
}
