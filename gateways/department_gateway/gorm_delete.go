package department_gateway

import (
	"conformity-core/frameworks/database/gorm/models"
)

func (g *GormDepartmentGatewayFacade) Delete(input DeleteInput) (bool, error) {
	result := g.DB.Conn.Where(
		"id = ?", input.DepartmentID,
	).Where(
		"company_id = ?", input.CompanyID,
	).Delete(&models.Department{})

	return result.RowsAffected > 0, result.Error
}
