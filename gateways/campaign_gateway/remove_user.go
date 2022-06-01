package campaign_gateway

import (
	"conformity-core/context"
	"conformity-core/frameworks/database/gorm/models"

	"gorm.io/gorm"
)

func (g GormCampaignGatewayFacade) RemoveUser(ctx *context.CoreCtx, input RemoveUserInput, tx *gorm.DB) (bool, error) {
	if tx == nil {
		tx = g.DB.Conn.WithContext(ctx.Context)
	}

	result := tx.Where(
		"department_id = ?",
		input.DepartmentID,
	).Where(
		"user_id = ?",
		input.UserID,
	).Delete(&models.UserCompanyDepartment{})

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
