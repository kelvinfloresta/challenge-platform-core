package import_case

import (
	database "conformity-core/frameworks/database/gorm"
	g "conformity-core/gateways/import_gateway"
	"conformity-core/usecases/user_case"
)

var Singleton = New(
	g.GormImportGatewayFacade{DB: database.DB_Production},
	&user_case.Singleton,
)

type ImportCase struct {
	gateway  g.IImportGateway
	userCase *user_case.UserCase
}

func New(
	gateway g.IImportGateway,
	userCase *user_case.UserCase,
) *ImportCase {
	return &ImportCase{gateway, userCase}
}
