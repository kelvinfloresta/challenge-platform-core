package user_case

import (
	"conformity-core/gateways/user_gateway"
	"conformity-core/usecases/campaign_case"
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/notification_case"
)

type UserCase struct {
	gateway          user_gateway.IUserGateway
	notificationCase *notification_case.NotificationCase
	companyCase      *company_case.CompanyCase
	campaignCase     *campaign_case.CampaignCase
}

func New(
	gateway user_gateway.IUserGateway,
	notificationCase *notification_case.NotificationCase,
	companyCase *company_case.CompanyCase,
	campaignCase *campaign_case.CampaignCase,
) UserCase {
	return UserCase{gateway, notificationCase, companyCase, campaignCase}
}

type LoginCaseInput struct {
	Email    string
	Password string
}

func (u UserCase) GetById(id string) (*user_gateway.GetByIdUserGatewayOutput, error) {
	return u.gateway.GetById(id)
}

func (u UserCase) GetByEmail(email string) (*user_gateway.GetByEmailUserGatewayOutput, error) {
	return u.gateway.GetByEmail(email)
}
