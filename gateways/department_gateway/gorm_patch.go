package department_gateway

import "conformity-core/frameworks/database/gorm/models"

func (g GormDepartmentGatewayFacade) Patch(input PatchInput) (bool, error) {
	result := g.DB.Conn.Model(&models.Department{}).Where(
		"company_id = ?", input.CompanyID,
	).Where(
		"id = ?", input.ID,
	).Update("name", input.Name)

	return result.RowsAffected > 0, result.Error
}
