package company_gateway

import (
	"conformity-core/enums/user_company_enum"
	"conformity-core/frameworks/database/gorm/models"
)

func (g GormCompanyGatewayFacade) ListManagers(input ListManagersInput) ([]ListManagersOutput, error) {

	query := g.DB.Conn.Model(&models.UserCompanyDepartment{}).Select(
		"email",
		"company_id",
	).Joins(
		"INNER JOIN users u ON u.id = user_id",
	).Joins(
		"INNER JOIN departments d ON d.id = department_id",
	).Where("role", user_company_enum.RoleCompanyManager)

	if len(input.CompanyIDs) > 0 {
		query = query.Where("d.company_id IN ?", input.CompanyIDs)
	}

	if input.Status != "" {
		query = query.Where("status = ?", input.Status)
	}

	output := []ListManagersOutput{}
	result := query.Scan(&output)

	return output, result.Error
}
