package department_gateway

import (
	"conformity-core/frameworks/database/gorm/models"
	"strings"
)

func (g *GormDepartmentGatewayFacade) GetOneByFilter(input GetOneByFilterInput) (*GetOneByFilterOutput, error) {
	data := &GetOneByFilterOutput{}

	result := g.DB.Conn.Model(&models.Department{}).Where(
		"company_id = ?", input.CompanyID,
	).Where(
		"lower(name) = ?", strings.ToLower(input.Name),
	).Limit(1).Scan(data)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return data, nil
}
