package campaign_case

import (
	c_gateway "conformity-core/gateways/campaign_gateway"
	"conformity-core/usecases/challenge_case"
	"conformity-core/usecases/company_case"
	"conformity-core/usecases/department_case"
	"conformity-core/usecases/notification_case"
)

type CampaignCase struct {
	gateway          c_gateway.ICampaignGateway
	departmentCase   department_case.DepartmentCase
	companyCase      company_case.CompanyCase
	challengeCase    *challenge_case.ChallengeCase
	notificationCase *notification_case.NotificationCase
}

func New(
	gateway c_gateway.ICampaignGateway,
	departmentCase department_case.DepartmentCase,
	companyCase company_case.CompanyCase,
	challengeCase *challenge_case.ChallengeCase,
	notificationCase *notification_case.NotificationCase,
) CampaignCase {
	return CampaignCase{gateway, departmentCase, companyCase, challengeCase, notificationCase}
}
