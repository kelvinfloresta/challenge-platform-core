package user_case

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/gateways/user_gateway"
	"conformity-core/usecases/campaign_case"
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/notification_case"
)

var Singleton UserCase = build()

func build() UserCase {
	return New(
		user_gateway.GormUserGatewayFacade{DB: database.DB_Production},
		notification_case.Singleton,
		&company_case.Singleton,
		&campaign_case.Singleton,
	)
}
