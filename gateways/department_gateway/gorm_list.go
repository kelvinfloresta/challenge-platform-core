package department_gateway

import (
	"conformity-core/frameworks/database/gorm/models"
)

func (g GormDepartmentGatewayFacade) List(input ListInput) ([]ListOutput, error) {
	query := g.DB.Conn.Model(&models.Department{}).Select("id, name")

	if input.CompanyID != "" {
		query = query.Where("company_id = ?", input.CompanyID)
	}

	output := []ListOutput{}
	result := query.Scan(&output)

	return output, result.Error
}
