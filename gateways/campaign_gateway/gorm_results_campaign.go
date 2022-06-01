package campaign_gateway

import (
	"conformity-core/context"
)

func (g GormCampaignGatewayFacade) ListResults(
	ctx *context.CoreCtx,
	input ListResultsInput,
) ([]ListResultsOutput, error) {
	data := []ListResultsOutput{}

	query := g.DB.Conn.WithContext(ctx.Context).Table("campaigns_results AS cr").Joins(
		`INNER JOIN users AS u ON u.id = cr.user_id`,
	).Joins(
		`INNER JOIN users_companies_departments AS ucd ON ucd.user_id = cr.user_id`,
	).Joins(
		`INNER JOIN challenges AS c ON c.id = cr.challenge_id`,
	).Joins(
		`INNER JOIN scheduled_challenges AS sc ON sc.challenge_id = cr.challenge_id AND sc.campaign_id = cr.campaign_id`,
	).Joins(
		`INNER JOIN departments AS d ON d.id = ucd.department_id`,
	).Select(`
		u.id as user_id,
		u.name as name,
		d.name as department_name,
		d.id as department_id,
		cr.tries,
		cr.correct,
		cr.challenge_id,
		c.title as challenge_title,
		sc.end_date,
		cr.campaign_id,
		cr.status
	`)

	if input.CampaignID != "" {
		query = query.Where(`cr.campaign_id = ?`, input.CampaignID)
	}

	if input.CompanyID != "" {
		query = query.Where(`d.company_id = ?`, input.CompanyID)
	}

	if input.ChallengeID != "" {
		query = query.Where(`c.id = ?`, input.ChallengeID)
	}

	if input.UserID != "" {
		query = query.Where(`cr.user_id = ?`, input.UserID)
	}

	if input.Status != "" {
		query = query.Where(`cr.status = ?`, input.Status)
	}

	result := query.Order("sc.end_date DESC, challenge_title ASC, name ASC").Scan(&data)
	return data, result.Error
}
