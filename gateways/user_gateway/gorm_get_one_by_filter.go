package user_gateway

import "conformity-core/frameworks/database/gorm/models"

func (g GormUserGatewayFacade) GetOneByFilter(input GetOneByFilterInput) (*GetOneByFilterOutput, error) {
	data := &GetOneByFilterOutput{}
	result := g.DB.Conn.Model(&models.User{}).Joins(`
		INNER JOIN users_companies_departments AS ucd
		ON ucd.user_id = users.id
	`).Where(
		"ucd.department_id = ?", input.DepartmentID,
	).Where(
		"ucd.user_id = ?", input.ID,
	).Select("users.phone").Scan(data)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return data, nil
}
