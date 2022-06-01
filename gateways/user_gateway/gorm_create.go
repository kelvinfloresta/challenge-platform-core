package user_gateway

import (
	"conformity-core/enums/user_company_enum"
	core_errors "conformity-core/errors"
	"conformity-core/frameworks/database/gorm/models"
	"conformity-core/utils"
	"strings"

	"gorm.io/gorm"
)

func AdaptError(err error) error {
	switch {

	case err == nil:
		return nil

	case strings.Contains(err.Error(), "unique_email"):
		return core_errors.ErrDuplicatedEmail

	case strings.Contains(err.Error(), "unique_login"):
		return core_errors.ErrDuplicatedLogin

	default:
		return err
	}

}

func (g GormUserGatewayFacade) Create(
	input CreateInput,
) (string, error) {
	userId, err := utils.UUID()
	if err != nil {
		return "", err
	}

	err = g.DB.Conn.Transaction(func(tx *gorm.DB) error {
		var email *string
		if input.Email != "" {
			email = &input.Email
		}

		err := tx.Create(&models.User{
			ID:          userId,
			Name:        input.Name,
			Password:    input.Password,
			Email:       email,
			Phone:       input.Phone,
			BirthDate:   input.BirthDate,
			Document:    input.Document,
			JobPosition: input.JobPosition,
		}).Error

		if err != nil {
			return err
		}

		err = tx.Create(&models.UserCompanyDepartment{
			UserID:       userId,
			OID:          input.OID,
			DepartmentID: input.DepartmentID,
			Login:        input.Login,
			Status:       user_company_enum.Active,
			Role:         input.Role,
		}).Error

		return err
	})

	return userId, AdaptError(err)
}
