package auth_case

import (
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/notification_case"
	"conformity-core/usecases/user_case"
)

type AuthCase struct {
	departmentCase   *department_case.DepartmentCase
	companyCase      *company_case.CompanyCase
	notificationCase *notification_case.NotificationCase
	userCase         *user_case.UserCase
}

func New(
	departmentCase *department_case.DepartmentCase,
	companyCase *company_case.CompanyCase,
	notificationCase *notification_case.NotificationCase,
	userCase *user_case.UserCase,
) *AuthCase {
	return &AuthCase{departmentCase, companyCase, notificationCase, userCase}
}

var Singleton = New(
	&department_case.Singleton,
	&company_case.Singleton,
	notification_case.Singleton,
	&user_case.Singleton,
)
