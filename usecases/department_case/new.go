package department_case

import g "conformity-core/gateways/department_gateway"

type DepartmentCase struct {
	gateway g.IDepartmentGateway
}

func New(gateway g.IDepartmentGateway) DepartmentCase {
	return DepartmentCase{gateway}
}
