package campaign_gateway

import (
	"conformity-core/context"

	"gorm.io/gorm/clause"
)

func (g GormCampaignGatewayFacade) GetAVGResult(
	ctx *context.CoreCtx,
	input GetAVGResultInput,
) ([]*GetAVGResultOutput, error) {
	selectFields := []string{
		`AVG(
			CASE
				WHEN correct = false
					THEN 0
				WHEN correct = true AND tries = 1
					THEN 100
				WHEN correct = true AND tries = 2
					THEN 67
				WHEN correct = true AND tries = 3
					THEN 34
				ELSE 0
			END
		)::int as AVG`,
	}

	query := g.DB.Conn.WithContext(ctx.Context).Table("campaigns_results AS cr").Joins(
		`INNER JOIN scheduled_challenges AS sc ON sc.challenge_id = cr.challenge_id
		AND cr.campaign_id = sc.campaign_id`,
	).Joins(
		`INNER JOIN departments AS d ON d.id = cr.department_id`,
	)

	if input.OR != nil {
		OR := make([]clause.Expression, 0)
		if input.OR.Correct != nil {
			OR = append(OR, clause.Or(clause.Eq{Column: "cr.correct", Value: *input.OR.Correct}))
		}
		if input.OR.Tries != nil {
			OR = append(OR, clause.Or(clause.Eq{Column: "cr.tries", Value: *input.OR.Tries}))
		}

		if !input.OR.EndDateLTE.IsZero() {
			OR = append(OR, clause.Or(clause.Lte{Column: "sc.end_date", Value: input.OR.EndDateLTE}))
		}

		and := clause.And(OR...)
		query.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{and},
		})
	}

	if input.Status != "" {
		query = query.Where("cr.status = ?", input.Status)
	}

	if input.CampaignID != "" {
		query = query.Where("cr.campaign_id = ?", input.CampaignID)
	}

	if input.CompanyID != "" {
		query = query.Where("d.company_id = ?", input.CompanyID)
	}

	if input.DepartmentID != "" {
		query = query.Where("cr.department_id = ?", input.DepartmentID)
	}

	if input.UserID != "" {
		query = query.Where("cr.user_id = ?", input.UserID)
	}

	if input.GroupBy == GroupByDepartment {
		selectFields = append(selectFields, "cr.department_id", "d.name as department_name")
		query = query.Group("department_id, d.name")
	}

	if input.GroupBy == GroupByUser {
		selectFields = append(selectFields, "cr.user_id")
		query = query.Group(GroupByUser)
	}

	output := []*GetAVGResultOutput{}
	result := query.Select(selectFields).Scan(&output)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return []*GetAVGResultOutput{}, nil
	}

	return output, nil
}
