package company_gateway

import "conformity-core/frameworks/database/gorm/models"

func (c GormCompanyGatewayFacade) GetOneByFilter(input GetOneByFilterInput) (*GetOneByFilterOutput, error) {
	data := &GetOneByFilterOutput{}
	query := c.DB.Conn.Model(&models.Company{})

	if input.DepartmentID != "" {
		query = query.Joins(`INNER JOIN departments d ON d.company_id = companies.id`).Where("d.id = ?", input.DepartmentID)
	}

	if input.Workspace != "" {
		query = query.Where("workspace = ?", input.Workspace)
	}

	if input.Domain != "" {
		query = query.Where("domain = ?", input.Domain)
	}

	result := query.Scan(data)

	if result.RowsAffected == 0 {
		return nil, result.Error
	}

	return data, result.Error
}
