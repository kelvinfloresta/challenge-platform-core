package department_case

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/gateways/department_gateway"
)

var Singleton DepartmentCase = build()

func build() DepartmentCase {
	return New(&department_gateway.GormDepartmentGatewayFacade{DB: database.DB_Production})
}
