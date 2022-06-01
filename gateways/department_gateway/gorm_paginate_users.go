package department_gateway

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"fmt"
	"strings"
)

func (g GormDepartmentGatewayFacade) PaginateUsers(input PaginateUsersInput) (*PaginateUsersOutput, error) {
	data := []PaginateUsers{}

	query := g.DB.Conn.Model(
		&models.User{},
	).Joins(
		"INNER JOIN users_companies_departments AS ucd ON ucd.user_id = users.id",
	).Joins(
		"INNER JOIN departments AS d ON ucd.department_id = d.id",
	).Where(
		"d.company_id = ?",
		input.CompanyID,
	).Where(
		"ucd.deleted_at IS NULL",
	)

	if input.Name != "" {
		query = query.Where("users.name ILIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(input.Name)))
	}

	if input.Email != "" {
		query = query.Where("users.email LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(input.Email)))
	}

	if input.DepartmentID != "" {
		query = query.Where("d.id = ?", input.DepartmentID)
	}

	if input.Status != "" {
		query = query.Where("ucd.status = ?", input.Status)
	}

	if input.Role != "" {
		query = query.Where("ucd.role = ?", input.Role)
	}

	var count int64
	users := query.Count(&count)
	if users.Error != nil {
		return nil, users.Error
	}

	if users.RowsAffected == 0 {
		return &PaginateUsersOutput{
			Data:     data,
			MaxPages: 0,
		}, nil
	}

	database.Paginate(query, database.PaginateInput{
		ActualPage: input.ActualPage,
		PageSize:   input.PageSize,
	})

	users = query.Select(
		"users.id",
		"users.name",
		"users.email",
		"users.document",
		"users.phone",
		"users.job_position",
		"users.created_at",
		"ucd.status",
		"ucd.role",
		"d.name as department_name",
		"d.id as department_id",
	).Scan(&data)

	if users.Error != nil {
		return nil, users.Error
	}

	output := &PaginateUsersOutput{
		Data:     data,
		MaxPages: database.CalcMaxPages(count, input.PageSize),
	}

	return output, nil
}
