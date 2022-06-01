package campaign_case

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/gateways/campaign_gateway"
	"conformity-core/usecases/challenge_case"
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/notification_case"
)

var Singleton CampaignCase = build()

func build() CampaignCase {
	gateway := campaign_gateway.GormCampaignGatewayFacade{DB: database.DB_Production}
	departmentCase := department_case.Singleton
	companyCase := company_case.Singleton
	challengeCase := challenge_case.Singleton
	notificationCase := notification_case.Singleton
	return New(gateway, departmentCase, companyCase, challengeCase, notificationCase)
}
