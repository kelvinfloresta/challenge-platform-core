package challenge_case

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/gateways/challenge_gateway"
)

var Singleton = build()

func build() *ChallengeCase {
	return New(&challenge_gateway.GormChallengeGatewayFacade{DB: database.DB_Production})
}
