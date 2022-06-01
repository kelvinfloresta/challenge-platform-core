package campaign_gateway

import (
	"conformity-core/context"
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormCampaignGatewayFacade struct {
	DB *database.Database
}

func (g GormCampaignGatewayFacade) Transaction(fn func(tx *gorm.DB) error) error {
	return g.DB.Conn.Transaction(fn)
}

func (g GormCampaignGatewayFacade) GetById(id string) (*GetByIdCampaignGatewayOutput, error) {
	data := &models.Campaign{ID: id}
	err := g.DB.Conn.Preload("Results").Preload("ScheduledChallenges").First(data).Error

	if err != nil {
		return nil, err
	}

	adapted := GetByIdCampaignGatewayOutput{
		ID:                  data.ID,
		Title:               data.Title,
		Results:             adaptResult(data.Results),
		ScheduledChallenges: adaptScheduledChallenges(data.ScheduledChallenges),
		CreatedAt:           data.CreatedAt,
		UpdatedAt:           data.UpdatedAt,
		DeletedAt:           data.DeletedAt.Time,
	}

	return &adapted, nil
}

func (g GormCampaignGatewayFacade) UpdateResult(input UpdateResultInput, tx *gorm.DB) (bool, error) {
	if tx == nil {
		tx = g.DB.Conn
	}

	if input.CampaignID != "" {
		tx = tx.Where("campaign_id = ?", input.CampaignID)
	}

	if input.ActualTries != nil {
		tx = tx.Where("tries = ?", input.ActualTries)
	}

	if input.ActualTriesLessThan != nil {
		tx = tx.Where("tries < ?", input.ActualTriesLessThan)
	}

	if input.ActualCorrect != nil {
		tx = tx.Where("correct = ?", input.ActualCorrect)
	}

	if input.ChallengeID != "" {
		tx = tx.Where("challenge_id = ?", input.ChallengeID)
	}

	if input.QuestionID != "" {
		tx = tx.Where("question_id = ?", input.QuestionID)
	}

	if input.UserID != "" {
		tx = tx.Where("user_id = ?", input.UserID)
	}

	if input.DepartmentID != "" {
		tx = tx.Where("department_id = ?", input.DepartmentID)
	}

	if input.EndDateGreaterThan != nil {
		tx = tx.Clauses(clause.From{
			Tables: []clause.Table{{Name: "scheduled_challenges"}},
		}).Where("end_date > ?", input.EndDateGreaterThan)
		tx.Statement.BuildClauses = append(tx.Statement.BuildClauses, "UPDATE", "SET", "FROM", "WHERE")
	}

	newValues := &models.CampaignResult{}
	if input.NewCorrect != nil {
		newValues.Correct = *input.NewCorrect
	}

	if input.NewTries != nil {
		newValues.Tries = *input.NewTries
	}

	if input.NewStatus != "" {
		newValues.Status = input.NewStatus
	}

	result := tx.Updates(newValues)

	if result.Error != nil {
		return false, result.Error
	}

	updated := result.RowsAffected > 0
	return updated, nil
}

func (g GormCampaignGatewayFacade) GetResult(ctx *context.CoreCtx, input GetResultInput) (*Result, error) {
	data := &models.CampaignResult{}
	result := g.DB.Conn.Limit(1).Find(data, input)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := &Result{
		UserID:       data.UserID,
		DepartmentID: data.DepartmentID,
		ChallengeID:  data.ChallengeID,
		CampaignID:   data.CampaignID,
		QuestionID:   data.QuestionID,
		Tries:        data.Tries,
		Correct:      data.Correct,
		Status:       data.Status,
	}

	return adapted, nil

}

func (g GormCampaignGatewayFacade) ListChallenges(ctx *context.CoreCtx, input ListChallengesInput) ([]GetChallengesCampaignOutput, error) {
	data := []GetChallengesCampaignOutput{}

	query := g.DB.Conn.Table("challenges AS c").Distinct(`
		c.id,
		c.title,
		c.media,
		c.questions,
		cr.campaign_id,
		sc.start_date,
		sc.end_date,
		AVG(
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
		)::int as impd
	`).Joins(`
		INNER JOIN scheduled_challenges AS sc
		ON c.id = sc.challenge_id

		INNER JOIN campaigns_results AS cr
		ON cr.campaign_id = sc.campaign_id
		AND cr.challenge_id = sc.challenge_id
	`).Group("c.id, cr.campaign_id, sc.start_date, sc.end_date")

	if !input.StartDateLTE.IsZero() {
		query = query.Where("sc.start_date <= ?", input.StartDateLTE)
	}

	if !input.EndDateGTE.IsZero() {
		query = query.Where("sc.end_date >= ?", input.EndDateGTE)
	}

	if input.UserID != "" {
		query = query.Where("cr.user_id = ?", input.UserID)
	}

	if input.Status != "" {
		query = query.Where("cr.status = ?", input.Status)
	}

	if input.TriesLT != nil {
		query = query.Where("cr.tries < ?", input.TriesLT)
	}

	if input.Correct != nil {
		query = query.Where("cr.correct = ?", input.Correct)
	}

	if input.CompanyID != "" {
		query = query.Joins(`
			INNER JOIN departments AS d
			ON d.id = cr.department_id
		`).Where("d.company_id = ?", input.CompanyID)
	}

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

	result := query.Scan(&data)

	return data, result.Error
}

type ListQuestionJSON struct {
}

func (g GormCampaignGatewayFacade) ListQuestions(ctx *context.CoreCtx, input ListQuestionsInput) ([]ListQuestionsOutput, error) {
	query := g.DB.Conn.Model(&models.CampaignResult{}).Preload("Challenge")

	if input.Status != "" {
		query = query.Where(`status = ?`, input.Status)
	}

	if input.Correct != nil {
		query = query.Where(`correct = ?`, input.Correct)
	}

	if input.TriesLessThan != nil {
		query = query.Where(`tries < ?`, input.TriesLessThan)
	}

	if input.UserID != "" {
		query = query.Where(`user_id = ?`, input.UserID)
	}

	if input.DepartmentID != "" {
		query = query.Where(`department_id = ?`, input.DepartmentID)
	}

	if input.ChallengeID != "" {
		query = query.Where(`challenge_id = ?`, input.ChallengeID)
	}

	if input.CampaignID != "" {
		query = query.Where(`campaign_id = ?`, input.CampaignID)
	}

	data := []models.CampaignResult{}
	result := query.Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	challengeQuestions := make(map[string]ListQuestionsOutput)
	for _, cr := range data {
		isEmpty := challengeQuestions[cr.QuestionID].Options == nil
		if !isEmpty {
			continue
		}

		questions := []ListQuestionsOutput{}
		if err := json.Unmarshal(cr.Challenge.Questions, &questions); err != nil {
			return nil, err
		}

		for _, lqo := range questions {
			challengeQuestions[lqo.ID] = lqo
		}

	}

	adapted := []ListQuestionsOutput{}
	for _, cr := range data {
		v := challengeQuestions[cr.QuestionID]
		v.Tries = cr.Tries
		adapted = append(adapted, v)
	}

	return adapted, nil
}

func (g GormCampaignGatewayFacade) List(ctx *context.CoreCtx, input ListInput) ([]*ListOutput, error) {
	data := []*ListOutput{}

	result := g.DB.Conn.Table("campaigns_results cr").Distinct(`
		cr.campaign_id,
		c.title,
		MIN(sc.start_date) AS start_date,
		MAX(sc.end_date) AS end_date
	`).Joins(`
		INNER JOIN campaigns c ON c.id = cr.campaign_id
		INNER JOIN departments d ON d.id = cr.department_id
		INNER JOIN scheduled_challenges sc ON sc.campaign_id = c.id
	`).Where(
		"d.company_id = ?", input.CompanyID,
	).Group(
		"cr.campaign_id, c.title",
	).Order(
		"end_date DESC",
	).Scan(&data)

	return data, result.Error
}

func (g GormCampaignGatewayFacade) ListUsers(ctx *context.CoreCtx, input ListUsersInput) ([]ListUsersOutput, error) {

	data := []ListUsersOutput{}
	query := g.DB.Conn.WithContext(ctx.Context).Table("users as u").Distinct(`
		u.email,
		u.name,
		u.id AS user_id,
		d.id AS department_id,
		d.company_id,
		d.name AS department_name,
		AVG(
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
		)::int as IMPD`,
	).Joins(
		`INNER JOIN users_companies_departments AS ucd ON u.id = ucd.user_id`,
	).Joins(
		`INNER JOIN campaigns_results AS cr ON cr.user_id = u.id`,
	).Joins(
		`INNER JOIN scheduled_challenges AS sc ON sc.campaign_id = cr.campaign_id`,
	).Joins(
		`INNER JOIN departments AS d ON d.id = ucd.department_id`,
	).Group("u.id, d.id")

	if input.UseCompanyStatus != nil {
		query = query.Where("ucd.status = ?", input.UseCompanyStatus)
	}

	if input.CampaignResultStatus != nil {
		query = query.Where("cr.status = ?", input.CampaignResultStatus)
	}

	if !input.StartDateLTE.IsZero() {
		query = query.Where("start_date <= ?", input.StartDateLTE)
	}

	if !input.StartDate.IsZero() {
		query = query.Where("start_date::DATE = ?", input.StartDate)
	}

	if !input.EndDateLTE.IsZero() {
		query = query.Where("end_date <= ?", input.EndDateLTE)
	}

	if !input.EndDateGT.IsZero() {
		query = query.Where("end_date > ?", input.EndDateGT)
	}

	if input.TriesLessThan != nil {
		query = query.Where("cr.tries < ?", input.TriesLessThan)
	}

	if input.Correct != nil {
		query = query.Where("cr.correct = ?", input.Correct)
	}

	if input.ChallengeID != "" {
		query = query.Where("cr.challenge_id = ?", input.ChallengeID)
	}

	if input.CampaignID != "" {
		query = query.Where("cr.campaign_id = ?", input.CampaignID)
	}

	result := query.Scan(&data)

	return data, result.Error
}

func (g GormCampaignGatewayFacade) ChangeUserStatus(input ChangeUserStatusInput, tx *gorm.DB) (bool, error) {
	if tx == nil {
		tx = g.DB.Conn
	}

	result := tx.Model(&models.UserCompanyDepartment{}).Where(
		"user_id = ? AND department_id = ?",
		input.UserID,
		input.DepartmentID,
	).Update("status", input.Status)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
