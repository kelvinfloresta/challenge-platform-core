package department_gateway

import (
	"conformity-core/frameworks/database/gorm/models"
)

func (g *GormDepartmentGatewayFacade) CountMonthlyUsers(input CountMonthlyUsersInput) ([]CountMonthlyUsersOutput, error) {
	data := make([]CountMonthlyUsersOutput, 0, 1000)

	result := g.DB.Conn.Model(
		&models.CampaignResult{},
	).Select(`
		COUNT(DISTINCT campaigns_results.user_id) AS count,
		c.id AS company_id,
		c.name AS company_name
	`).Joins(`
		INNER JOIN departments AS d ON d.id = campaigns_results.department_id
		INNER JOIN companies AS c ON c.id = d.company_id
		INNER JOIN scheduled_challenges AS sc
			ON  sc.campaign_id  = campaigns_results.campaign_id
			AND sc.challenge_id = campaigns_results.challenge_id
	`).Where(
		"sc.start_date >= ?", input.StartDateGTE,
	).Where(
		"sc.end_date <= ?", input.EndDateLTE,
	).Group(
		"c.id, c.name",
	).Scan(&data)

	return data, result.Error
}
