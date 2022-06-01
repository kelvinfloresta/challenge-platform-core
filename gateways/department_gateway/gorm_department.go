package department_gateway

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"conformity-core/utils"
)

type GormDepartmentGatewayFacade struct {
	DB *database.Database
}

func (g *GormDepartmentGatewayFacade) Create(input CreateDepartmentGatewayInput) (string, error) {
	id, err := utils.UUID()
	if err != nil {
		return "", err
	}

	result := g.DB.Conn.Create(&models.Department{
		ID:        id,
		Name:      input.Name,
		CompanyID: input.CompanyID,
	})

	if result.Error != nil {
		return "", result.Error
	}

	return id, nil
}

func (g GormDepartmentGatewayFacade) GetById(id string) (*GetByIdDepartmentGatewayOutput, error) {
	data := &models.Department{}
	result := g.DB.Conn.Find(data, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := &GetByIdDepartmentGatewayOutput{
		ID:        data.ID,
		Name:      data.Name,
		CompanyID: data.CompanyID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.Time,
	}

	return adapted, nil
}

func (g GormDepartmentGatewayFacade) GetUsersByFilter(
	input GetUsersByFilterGatewayInput,
) (out *[]GetUsersByFilterGatewayOutput, err error) {
	out = &[]GetUsersByFilterGatewayOutput{}

	query := g.DB.Conn.Model(
		&models.UserCompanyDepartment{},
	).Select(
		"users.name",
		"users.email",
		"users.id AS user_id",
		"departments.id AS department_id",
	).Joins(
		`INNER JOIN departments ON departments.id = users_companies_departments.department_id`,
	).Joins(
		`INNER JOIN companies ON departments.company_id = companies.id`,
	).Joins(
		`INNER JOIN users ON users_companies_departments.user_id = users.id`,
	)

	if !input.UserCreatedAtGTE.IsZero() {
		query = query.Where("users_companies_departments.created_at >= ?", input.UserCreatedAtGTE)
	}

	if !input.UserCreatedAtLTE.IsZero() {
		query = query.Where("users_companies_departments.created_at <= ?", input.UserCreatedAtLTE)
	}

	if input.UserCompanyStatus != "" {
		query = query.Where("status", input.UserCompanyStatus)
	}

	if input.CompanyID != "" {
		query.Where("companies.id", input.CompanyID)
	}

	if len(input.Departments) > 0 {
		query.Where("departments.id", input.Departments)
	}

	err = query.Find(out).Error

	return
}

func (g GormDepartmentGatewayFacade) GetUserByLogin(login string) (*UserDepartment, error) {
	data := &UserDepartment{}
	result := g.DB.Conn.Model(&models.UserCompanyDepartment{}).Select(
		"users.name",
		"users.password",
		"companies.id as company_id",
		"companies.require_password",
		"companies.sso_enabled",
		"companies.workspace",
		"users_companies_departments.login",
		"users_companies_departments.role",
		"users_companies_departments.user_id",
		"users_companies_departments.o_id",
		"users_companies_departments.department_id",
		"users_companies_departments.status",
		"users_companies_departments.created_at",
		"users_companies_departments.updated_at",
		"users_companies_departments.deleted_at",
	).Joins(
		`INNER JOIN departments ON departments.id = users_companies_departments.department_id`,
	).Joins(
		"INNER JOIN users ON users.id = users_companies_departments.user_id",
	).Joins(
		"INNER JOIN companies ON departments.company_id = companies.id",
	).Where(
		"login = ?", login,
	).Scan(data)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return data, nil
}
