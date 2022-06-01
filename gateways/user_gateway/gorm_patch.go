package user_gateway

import (
	"conformity-core/frameworks/database/gorm/models"

	"gorm.io/gorm"
)

func (g GormUserGatewayFacade) Patch(input PatchInput) (bool, error) {
	user := models.User{}
	newValues := map[string]interface{}{
		"name":         input.Name,
		"job_position": input.JobPosition,
		"phone":        input.Phone,
		"email":        &input.Email,
		"document":     input.Document,
	}

	updated := false
	err := g.DB.Conn.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&user).Where(
			"id = ?", input.ID,
		).Updates(newValues)

		if result.Error != nil {
			return result.Error
		}

		updated = result.RowsAffected > 0
		if !updated {
			return nil
		}

		if input.OID == "" && input.Role == "" && input.Login == "" {
			return nil
		}

		userCompany := models.UserCompanyDepartment{
			OID:   input.OID,
			Role:  input.Role,
			Login: input.Login,
		}

		return tx.Table(userCompany.TableName()).Where(
			"user_id = ?", input.ID,
		).Where(
			"department_id = ?", input.DepartmentID,
		).Updates(userCompany).Error
	})

	return updated, AdaptError(err)
}
