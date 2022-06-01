package import_gateway

import (
	"conformity-core/enums/user_company_enum"
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"conformity-core/gateways/user_gateway"
	"conformity-core/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type GormImportGatewayFacade struct {
	DB *database.Database
}

func (g GormImportGatewayFacade) getDepartmentId(companyId, departmentName string) (string, error) {
	id := ""

	result := g.DB.Conn.Model(&models.Department{}).Where(
		"company_id = ?", companyId,
	).Where(
		"lower(name) = ?", strings.ToLower(departmentName),
	).Limit(1).Select("id").Scan(&id)

	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", nil
	}

	return id, nil
}

func (g GormImportGatewayFacade) ImportUsers(input *ImportUsersInput) ([]ImportUsersOutput, error) {
	uniqueDepartments := make(map[string]string)
	departmentsAlreadyCreated := make(map[string]string)
	users := make([]models.User, 0, len(input.Data))
	usersDepartments := []models.UserCompanyDepartment{}
	importedUsers := make([]ImportUsersOutput, 0, len(input.Data))

	for _, iud := range input.Data {
		departmentName := strings.Trim(iud.DepartmentName, " ")
		email := strings.ToLower(strings.Trim(iud.Email, " "))
		userName := strings.Trim(iud.UserName, " ")
		login := email

		if departmentsAlreadyCreated[departmentName] == "" &&
			uniqueDepartments[departmentName] == "" {
			departmentId, err := g.getDepartmentId(input.CompanyID, departmentName)
			if err != nil {
				return nil, err
			}

			if departmentId != "" {
				departmentsAlreadyCreated[departmentName] = departmentId
			} else {
				departmentId, err := utils.UUID()
				if err != nil {
					return nil, err
				}
				uniqueDepartments[departmentName] = departmentId
			}
		}

		userId, err := utils.UUID()
		if err != nil {
			return nil, err
		}

		users = append(users, models.User{
			ID:    userId,
			Name:  userName,
			Email: &email,
		})

		departmentId := uniqueDepartments[departmentName]
		if departmentId == "" {
			departmentId = departmentsAlreadyCreated[departmentName]
		}

		usersDepartments = append(usersDepartments, models.UserCompanyDepartment{
			Login:        login,
			Status:       user_company_enum.Active,
			Role:         user_company_enum.RoleUser,
			UserID:       userId,
			DepartmentID: departmentId,
		})

		importedUsers = append(importedUsers, ImportUsersOutput{
			ID:           userId,
			Name:         userName,
			Email:        email,
			Login:        login,
			DepartmentID: departmentId,
		})

	}

	departments := make([]models.Department, 0, len(uniqueDepartments))
	for departmentName, departmentId := range uniqueDepartments {
		departments = append(departments, models.Department{
			ID:        departmentId,
			Name:      departmentName,
			CompanyID: input.CompanyID,
		})
	}

	err := g.DB.Conn.Transaction(func(tx *gorm.DB) error {
		if len(departments) != 0 {
			err := tx.Create(departments).Error
			if err != nil {
				return fmt.Errorf("create department: %v", err)
			}
		}

		err := tx.Create(users).Error
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}

		err = tx.Create(usersDepartments).Error
		if err != nil {
			return fmt.Errorf("create users departments: %v", err)
		}

		return nil

	})

	return importedUsers, user_gateway.AdaptError(err)
}
