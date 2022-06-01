package company_case

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/gateways/company_gateway"
)

var Singleton CompanyCase = build()

func build() CompanyCase {
	gateway := &company_gateway.GormCompanyGatewayFacade{DB: database.DB_Production}
	return New(gateway)
}
