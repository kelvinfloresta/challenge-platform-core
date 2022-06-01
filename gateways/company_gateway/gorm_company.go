package company_gateway

import (
	"conformity-core/context"
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"conformity-core/utils"

	"gorm.io/gorm"
)

type GormCompanyGatewayFacade struct {
	DB *database.Database
}

func (g *GormCompanyGatewayFacade) GetTransaction(ctx *context.CoreCtx) *gorm.DB {
	return g.DB.Conn.Begin().WithContext(ctx.Context)
}

func (g *GormCompanyGatewayFacade) Create(input CreateCompanyGatewayInput) (string, error) {
	id, err := utils.UUID()
	if err != nil {
		return "", err
	}

	result := g.DB.Conn.Create(&models.Company{
		ID:              id,
		Name:            input.Name,
		Document:        input.Document,
		Workspace:       input.Workspace,
		Domain:          input.Domain,
		RequirePassword: input.RequirePassword,
	})

	if result.Error != nil {
		return "", result.Error
	}

	return id, nil
}

func (g GormCompanyGatewayFacade) GetById(id string) (*GetByIdOutput, error) {
	data := &models.Company{}
	result := g.DB.Conn.Limit(1).Find(data, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIdOutput{
		ID:              data.ID,
		Name:            data.Name,
		Document:        data.Document,
		Workspace:       data.Workspace,
		RequirePassword: data.RequirePassword,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
		DeletedAt:       data.DeletedAt.Time,
	}

	return &adapted, nil
}

func (g GormCompanyGatewayFacade) GetUser(input GetUserInput) (*GetUserOutput, error) {
	data := &GetUserOutput{}

	result := g.DB.Conn.Limit(1).Model(
		&models.UserCompanyDepartment{},
	).Where(
		"user_id = ?", input.UserID,
	).Where(
		"department_id = ?", input.DepartmentID,
	).Scan(data)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return data, nil
}

func (g GormCompanyGatewayFacade) ChangeUserRole(input ChangeUserRoleInput, tx *gorm.DB) (bool, error) {
	if tx == nil {
		tx = g.DB.Conn
	}

	result := tx.Model(&models.UserCompanyDepartment{}).Where(
		"user_id = ? AND department_id = ?",
		input.UserID,
		input.DepartmentID,
	).Update("role", input.Role)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
