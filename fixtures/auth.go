package fixtures

import (
	"conformity-core/usecases/auth_case"
)

type authCaseFixture struct{}

var AuthCase authCaseFixture = authCaseFixture{}

func (u authCaseFixture) NewIntegration() *auth_case.AuthCase {
	departmentCase := DepartmentCase.NewIntegration()
	companyCase := CompanyCase.NewIntegration()
	notificationCase, _ := NotificationCase.NewUnit()
	userCase := UserCase.NewIntegration()
	return auth_case.New(&departmentCase, &companyCase, notificationCase, &userCase)

}

// #nosec
const TestToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
